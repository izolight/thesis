package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import Signature
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SigningRequest
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.*
import com.auth0.jwt.interfaces.DecodedJWT
import com.google.protobuf.ByteString
import io.ktor.application.call
import io.ktor.request.receive
import io.ktor.routing.Routing
import io.ktor.routing.post
import org.koin.ktor.ext.inject
import java.io.File

fun Routing.sign() {
    val oidcService by inject<IOIDCService>()
    val secretService by inject<ISecretService>()
    val signingKeyService by inject<ISigningKeysService>()
    val caService by inject<ICertificateAuthorityService>()
    val tsaService by inject<ITimestampingService>()

    fun validateOidcNonce(idToken: DecodedJWT, salt: ByteArray, concatenatedHashes: ByteArray) {
        val oidcNonce = hmacSha256(salt, concatenatedHashes)
        val oidcNonceAsHexString = byteArrayToHexString(oidcNonce)
        if (idToken.getClaim("nonce").asString() != oidcNonceAsHexString) {
            throw InvalidDataException(
                "Nonce mismatch"
            )
        }
    }

    fun validateSalt(signingRequest: Valid<SigningRequest>): ByteArray {
        val hmacKey = secretService.getSecret()
        val concatenatedHashes = signingRequest.value.hashes.joinToString("").toByteArray()
        val salt = hmacSha256(hmacKey, concatenatedHashes)
        val saltAsHexString = byteArrayToHexString(salt)

        if (saltAsHexString != signingRequest.value.salt) {
            throw InvalidDataException(
                "Salt mismatch"
            )
        } else {
            return salt
        }
    }

    post(URLs.SIGN) {
        when (val input = call.receive<SigningRequest>().validate()) {
            is Valid -> {
                val salt = validateSalt(input)
                val jwtValidationResult = oidcService.validateIdToken(input.value.id_token)
                validateOidcNonce(
                    jwtValidationResult.idToken,
                    salt,
                    input.value.hashes.joinToString("").toByteArray(Charsets.UTF_8)
                )
                when (val subjectInformation = SigningKeySubjectInformation.fromIdToken(jwtValidationResult.idToken)) {
                    is Either.Success -> {
                        caService.signCSR(
                            signingKeyService.generateSigningKey(subjectInformation.value)
                        ).also {
                            signingKeyService.signToPkcs7(
                                subjectInformation.value,
                                Signature.SignatureData.newBuilder()
                                    .addAllSaltedDocumentHash(
                                        input.value.hashes.map { s ->
                                            hmacSha256(salt, hexStringToByteArray(s)).toByteString()
                                        }
                                    )
                                    .setHashAlgorithm(Signature.HashAlgorithm.SHA2_256)
                                    .setMacKey(salt.toByteString())
                                    .setMacAlgorithm(Signature.MACAlgorithm.HMAC_SHA2_256)
                                    .setSignatureLevel(Signature.SignatureLevel.ADVANCED)
                                    .setIdToken(ByteString.copyFromUtf8(input.value.id_token))
                                    .setJwkIdp(
                                        ByteString.copyFromUtf8(
                                            oidcService.marshalJwk(jwtValidationResult.jwk)
                                        )
                                    )
                                    .build().toByteArray(),
                                it,
                                caService.fetchBundleAsync(it)
                            ).encoded.also { pkcs7Signature ->
                                Signature.SignatureFile.newBuilder()
                                    .setSignatureDataInPkcs7(pkcs7Signature.toByteString())
                                    .addRfc3161InPkcs7(tsaService.stamp(pkcs7Signature).toByteString())
                                    .build()
                            }.also { signatureFile ->
                                File("/tmp/signaturefile").writeBytes(signatureFile)
                            }
                        }
                        println()
                    }
                    is Either.Error -> throw subjectInformation.e
                }
            }
            is Invalid -> throw input.error
        }
    }
}

