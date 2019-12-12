package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.HashesSubmissionResponse
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.Config
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.INonceGeneratorService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IOIDCService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.ISecretService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.*
import io.ktor.application.call
import io.ktor.request.receive
import io.ktor.response.respond
import io.ktor.routing.Routing
import io.ktor.routing.post
import org.koin.ktor.ext.inject
import org.slf4j.Logger

fun Routing.postHashes() {
    val nonceGenerator by inject<INonceGeneratorService>()
    val secretService by inject<ISecretService>()
    val oidcService by inject<IOIDCService>()
    val logger by inject<Logger>()

    post(URLs.SUBMIT_HASHES) {
        when (val input = call.receive<SubmittedHashes>().validate()) {
            is Valid -> {
                val seed = nonceGenerator.getNonce()
                val sortedHashes = input.value.hashes.sorted()
                val salt = calculateSalt(secretService.hkdf(seed), sortedHashes)
                val oidcNonce = calculateOidcNonce(maskHashes(sortedHashes, salt).concatenate())

                val idpRedirect = oidcService.constructAuthenticationRequestUrl(
                    oidcService.getAuthorisationEndpoint(),
                    nonce = oidcNonce,
                    state = oidcNonce
                )

                val response = HashesSubmissionResponse(
                        providers = mapOf(Config.OIDC_IDP_NAME to idpRedirect.toString()),
                        salt = byteArrayToHexString(salt),
                        seed = byteArrayToHexString(seed)
                    )

                call.respond(response)

                logger.info("Received hashes: ${input.value.hashes}")
                logger.info("Calculated salt ${byteArrayToHexString(salt)} and seed ${byteArrayToHexString(seed)}")
                logger.info("Responding: ${response}")
            }
            is Invalid -> throw input.error
        }
    }
}

