package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import Signature
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SigningRequest
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.Either
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.byteArrayToHexString
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hexStringToByteArray
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import com.auth0.jwt.interfaces.DecodedJWT
import com.google.protobuf.ByteString
import io.ktor.application.call
import io.ktor.request.receive
import io.ktor.routing.Routing
import io.ktor.routing.post
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration
import org.koin.ktor.ext.inject
import org.slf4j.LoggerFactory
import java.io.File

fun Routing.sign() {
    val oidcService by inject<IOIDCService>()
    val secretService by inject<ISecretService>()
    val signingKeyService by inject<ISigningKeysService>()
    val caService by inject<ICertificateAuthorityService>()
    val tsaService by inject<ITimestampingService>()
    val logger = LoggerFactory.getLogger(this.javaClass)
    val json = Json(JsonConfiguration.Stable)

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
                        val cert = caService.signCSR(
                            signingKeyService.generateSigningKey(subjectInformation.value)
                        )
                        val bundle = caService.fetchBundle(cert)

//                        TODO("other hashes, macced, sorted")
                        input.value.hashes.map { s -> ByteString.copyFrom(hmacSha256(salt, hexStringToByteArray(s))) }
                        val signatureData = Signature.SignatureData.newBuilder()
                            .addAllSaltedDocumentHash(
                                input.value.hashes.map { s ->
                                    ByteString.copyFrom(
                                        hmacSha256(salt, hexStringToByteArray(s))
                                    )
                                }
                            )
                            .setHashAlgorithm(Signature.HashAlgorithm.SHA2_256)
                            .setMacKey(ByteString.copyFrom(salt))
                            .setMacAlgorithm(Signature.MACAlgorithm.HMAC_SHA2_256)
                            .setSignatureLevel(Signature.SignatureLevel.ADVANCED)
                            .setIdToken(ByteString.copyFromUtf8(input.value.id_token))
                            .addJwkIdp(
                                ByteString.copyFromUtf8(
                                    oidcService.marshalJwk(jwtValidationResult.jwk)
                                )
                            )
                            .build()

                        val pkcs7 = signingKeyService.signToPkcs7(
                            subjectInformation.value,
                            signatureData.toByteArray(),
                            cert,
                            bundle
                        ).encoded

                        val timestampOfSignatureData = tsaService.stamp(pkcs7)
                        val signatureFile = Signature.SignatureFile.newBuilder()
                            .setSignatureDataInPkcs7(ByteString.copyFrom(pkcs7))
                                // TODO add pkcs7 enveloped timestamp
                            .build()

                        File("/tmp/signaturefile").writeBytes(signatureFile.toByteArray())
                        println()
                    }
                    is Either.Error -> throw subjectInformation.e
                }

            }
            is Invalid -> throw input.error
        }

    }


}

