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
import kotlinx.serialization.UnstableDefault
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration
import org.koin.ktor.ext.inject
import org.slf4j.Logger

@UnstableDefault
fun Routing.postHashes() {
    val nonceGenerator by inject<INonceGeneratorService>()
    val secretService by inject<ISecretService>()
    val oidcService by inject<IOIDCService>()
    val logger by inject<Logger>()
    val prettyJson = Json(JsonConfiguration.Default.copy(prettyPrint = true))

    post(URLs.SUBMIT_HASHES) {
        when (val input = call.receive<SubmittedHashes>().validate()) {
            is Valid -> {
                logger.info("Request: {}", prettyJson.stringify(SubmittedHashes.serializer(), input.value))
                val seed = nonceGenerator.getNonce()
                val sortedHashes = input.value.hashes.sorted()
                val salt = calculateSalt(secretService.hkdf(seed), sortedHashes)
                val oidcNonce = calculateOidcNonce(maskHashes(sortedHashes, salt).concatenate())
                logger.info("Generated seed: {}", byteArrayToHexString(seed))
                logger.info("Calculated salt: {}", byteArrayToHexString(salt))
                logger.info("Calculated OIDC nonce: {}", oidcNonce)

                logger.info("Fetching OIDC configuration and constructing redirect URIs")
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

                logger.info("Response: {}", prettyJson.stringify(HashesSubmissionResponse.serializer(), response))
            }
            is Invalid -> throw input.error
        }
    }
}

