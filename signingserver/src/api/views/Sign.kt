package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import Signature
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.*
import com.google.protobuf.ByteString
import io.ktor.application.call
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.locations
import io.ktor.request.receive
import io.ktor.response.respond
import io.ktor.routing.Routing
import io.ktor.routing.post
import kotlinx.serialization.UnstableDefault
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration
import org.bouncycastle.asn1.ASN1Encoding
import org.bouncycastle.asn1.ASN1OutputStream
import org.bouncycastle.cms.CMSSignedData
import org.koin.ktor.ext.inject
import org.slf4j.Logger
import java.io.ByteArrayOutputStream

@UnstableDefault
@KtorExperimentalLocationsAPI
fun Routing.sign() {
    val oidcService by inject<IOIDCService>()
    val secretService by inject<ISecretService>()
    val signingKeyService by inject<ISigningKeysService>()
    val caService by inject<ICertificateAuthorityService>()
    val tsaService by inject<ITimestampingService>()
    val signatureHoldingService by inject<ISignaturesHoldingService>()
    val logger by inject<Logger>()
    val prettyJson = Json(JsonConfiguration.Default.copy(prettyPrint = true))

    // TODO refactor this monster into smaller methods
    post(URLs.SIGN) {
        when (val input = call.receive<SigningRequest>().validate()) {
            is Valid -> {
                logger.info("Request: {}", prettyJson.stringify(SigningRequest.serializer(), input.value))
                val sortedHashes = input.value.hashes.sorted()
                val salt = calculateSalt(
                    secretService.hkdf(hexStringToByteArray(input.value.seed)),
                    sortedHashes
                )
                if (byteArrayToHexString(salt) != input.value.salt) {
                    throw InvalidDataException("Salt mismatch")
                } else {
                    logger.info("Salt validation succeeded")
                }

                val jwtValidationResult = oidcService.validateIdToken(input.value.id_token)
                val maskedHashes = maskHashes(sortedHashes, salt)
                if (calculateOidcNonce(maskedHashes.concatenate()) != jwtValidationResult.idToken.getClaim("nonce").asString()) {
                    throw InvalidDataException("Nonce mismatch")
                }
                logger.info("Seed, salt, nonce, id_token validation succeeded")

                val subjectInformation = SigningKeySubjectInformation.fromIdToken(jwtValidationResult.idToken)
                logger.info(
                    "Generating signing key for subject {}",
                    prettyJson.stringify(SigningKeySubjectInformation.serializer(), subjectInformation)
                )
                val certificateHolder = caService.signCSR(
                    signingKeyService.generateSigningKey(subjectInformation).also {
                        logger.info("Requesting CA to sign signing key")
                    }
                )
                logger.info("Constructing and signing inner CMS")
                val pkcs7Signature = signingKeyService.signToPkcs7(
                    subjectInformation,
                    buildSignature(
                        maskedHashes,
                        salt,
                        input.value.id_token,
                        oidcService.marshalJwk(jwtValidationResult.jwk)
                    ),
                    certificateHolder
                ).toDER()

                logger.info(
                    "Destroying signing key for subject {}", prettyJson.stringify(
                        SigningKeySubjectInformation.serializer(),
                        subjectInformation
                    )
                )
                signingKeyService.destroySigningKey(subjectInformation)
                logger.info("Requesting timestamp from TSA")
                val signatureFile = buildSignaturefile(
                    pkcs7Signature,
                    timestamp = tsaService.stamp(pkcs7Signature)
                )
                signatureHoldingService.generateId().also { id ->
                    signatureHoldingService.set(id, signatureFile.toByteArray())
                    logger.info("Signing successful, storing signature with ID {}", id)
                    call.respond(
                        SigningResponse(
                            signature =
                            "${URLs.BASE_URL}/${locations.href(SignatureRetrievalRequest(id = id))}"
                        )
                    )
                }
            }
            is Invalid -> throw input.error
        }
    }
}

fun buildSignature(
    maskedHashes: List<ByteArray>,
    salt: ByteArray,
    idToken: String,
    jwk: String
): Signature.SignatureData = Signature.SignatureData.newBuilder()
    .addAllSaltedDocumentHash(
        maskedHashes.map { it.toByteString() }
    )
    .setHashAlgorithm(Signature.HashAlgorithm.SHA2_256)
    .setMacKey(salt.toByteString())
    .setMacAlgorithm(Signature.MACAlgorithm.HMAC_SHA2_256)
    .setSignatureLevel(Signature.SignatureLevel.ADVANCED)
    .setIdToken(ByteString.copyFromUtf8(idToken))
    .setJwkIdp(ByteString.copyFromUtf8(jwk))
    .build()

fun buildSignaturefile(
    pkcs7Signature: ByteArray,
    timestamp: ByteArray
): Signature.SignatureFile = Signature.SignatureFile.newBuilder()
    .setSignatureData(pkcs7Signature.toByteString())
    .addRfc3161(timestamp.toByteString())
    .build()

fun CMSSignedData.toDER(): ByteArray = ByteArrayOutputStream().also {
    ASN1OutputStream.create(it, ASN1Encoding.DER).also { asn1outputStream ->
        asn1outputStream.writeObject(this.toASN1Structure())
        asn1outputStream.close()
    }
}.toByteArray()
