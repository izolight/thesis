package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.HashesSubmissionResponse
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def.ISecretService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl.Config
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl.INonceGeneratorService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl.IOIDCService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.*
import io.ktor.application.*
import io.ktor.request.*
import io.ktor.response.*
import io.ktor.routing.*
import kotlinx.serialization.json.Json
import org.koin.ktor.ext.inject
import org.slf4j.Logger

fun Routing.postHashes() {
    val nonceGenerator by inject<INonceGeneratorService>()
    val secretService by inject<ISecretService>()
    val oidcService by inject<IOIDCService>()
    val logger by inject<Logger>()
    val prettyJson = Json { prettyPrint = true }

    post(URLs.SUBMIT_HASHES) {
        when (val input = call.receive<SubmittedHashes>().validate()) {
            is Valid -> {
                logger.info("Request: {}", prettyJson.encodeToString(SubmittedHashes.serializer(), input.value))
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

                logger.info("Response: {}", prettyJson.encodeToString(HashesSubmissionResponse.serializer(), response))
            }
            is Invalid -> throw input.error
        }
    }
}

