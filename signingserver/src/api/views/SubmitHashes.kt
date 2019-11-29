package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.HashesSubmissionResponse
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.Config
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.INonceGeneratorService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IOIDCService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.ISecretService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.byteArrayToHexString
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import io.ktor.application.call
import io.ktor.request.receive
import io.ktor.response.respond
import io.ktor.routing.Routing
import io.ktor.routing.post
import org.koin.ktor.ext.inject
import org.slf4j.LoggerFactory

fun Routing.postHashes() {
    val logger = LoggerFactory.getLogger(this.javaClass)
    val nonceGenerator by inject<INonceGeneratorService>()
    val secretService by inject<ISecretService>()
    val oidcService by inject<IOIDCService>()

    post(URLs.SUBMIT_HASHES) {
        when (val input = call.receive<SubmittedHashes>().validate()) {
            is Valid -> {
                val seed = nonceGenerator.getNonce()
                val seedAsHexString = byteArrayToHexString(seed)
                val hmacKey = secretService.getSecret()
                val hmacKeyAsHexString = byteArrayToHexString(hmacKey)
                val concatenatedHashes = input.value.hashes.joinToString("")
                val concatenatedHashesAsByteArray = concatenatedHashes.toByteArray(Charsets.UTF_8)
                val salt = hmacSha256(hmacKey, concatenatedHashesAsByteArray)
                val saltAsHexString = byteArrayToHexString(salt)
                val oidcNonce = hmacSha256(salt, concatenatedHashesAsByteArray)
                val oidcNonceAsHexString = byteArrayToHexString(oidcNonce)
                logger.debug("hmacKey: {}", hmacKeyAsHexString)
                val idpRedirect = oidcService.constructAuthenticationRequestUrl(
                    oidcService.getAuthorisationEndpoint(),
                    nonce = oidcNonceAsHexString,
                    state = oidcNonceAsHexString
                )
                call.respond(
                    HashesSubmissionResponse(
                        mapOf(Config.OIDC_IDP_NAME to idpRedirect.toString()),
                        salt = saltAsHexString,
                        seed = seedAsHexString
                    )
                )
            }
            is Invalid -> throw input.error
        }
    }
}

