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
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hexStringToByteArray
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import io.ktor.application.call
import io.ktor.request.receive
import io.ktor.response.respond
import io.ktor.routing.Routing
import io.ktor.routing.post
import org.koin.ktor.ext.inject
import org.slf4j.LoggerFactory
import java.io.ByteArrayOutputStream

fun Routing.postHashes() {
    val logger = LoggerFactory.getLogger(this.javaClass)
    val nonceGenerator by inject<INonceGeneratorService>()
    val secretService by inject<ISecretService>()
    val oidcService by inject<IOIDCService>()

    post(URLs.SUBMIT_HASHES) {
        when (val input = call.receive<SubmittedHashes>().validate()) {
            is Valid -> {
                val seed = nonceGenerator.getNonce()

                val hmacKey = secretService.hkdf(seed)

                val sortedHashes = input.value.hashes.sorted()

                val salt = hmacSha256(hmacKey, hexStringToByteArray(sortedHashes.joinToString("")))

                val maskedHashes = ByteArrayOutputStream(
                    // pre-allocate buffer
                    sortedHashes.sumBy { it.length }
                ).also {
                    sortedHashes.map { h ->
                        it.write(
                            hmacSha256(salt, hexStringToByteArray(h))
                        )
                    }
                }.toByteArray()

                val oidcNonce = byteArrayToHexString(hmacSha256(salt, maskedHashes))

                val idpRedirect = oidcService.constructAuthenticationRequestUrl(
                    oidcService.getAuthorisationEndpoint(),
                    nonce = oidcNonce,
                    state = oidcNonce
                )

                call.respond(
                    HashesSubmissionResponse(
                        providers = mapOf(Config.OIDC_IDP_NAME to idpRedirect.toString()),
                        salt = byteArrayToHexString(salt),
                        seed = byteArrayToHexString(seed)
                    )
                )
            }
            is Invalid -> throw input.error
        }
    }
}

