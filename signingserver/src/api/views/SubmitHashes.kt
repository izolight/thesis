package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.HashesSubmissionResponse
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.INonceGeneratorService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IOIDCService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.ISecretService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.byteArrayToHexString
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hexStringToByteArray
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import io.ktor.application.call
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.Location
import io.ktor.locations.post
import io.ktor.request.receive
import io.ktor.response.respond
import io.ktor.routing.Routing
import org.koin.ktor.ext.inject

@KtorExperimentalLocationsAPI
@Location("/api/v1/hashes")
class HashesRoute()

@KtorExperimentalLocationsAPI
fun Routing.postHashes() {
    val nonceGenerator by inject<INonceGeneratorService>()
    val secretService by inject<ISecretService>()
    val oidcService by inject<IOIDCService>()

    post<HashesRoute> {
        when (val input = call.receive<SubmittedHashes>().validate()) {
            is Valid -> {
                val seed = nonceGenerator.getNonce()
                val hmacKey = secretService.getSecret()
                val concatenatedHashes = hexStringToByteArray(input.value.hashes.joinToString(""))
                val salt = hmacSha256(hmacKey, concatenatedHashes)
                val oidcNonce = hmacSha256(salt, concatenatedHashes)
                val oidcNonceAsHexString = byteArrayToHexString(oidcNonce)
                val idpRedirect = oidcService.constructAuthenticationRequestUrl(
                    oidcService.getAuthorisationEndpoint(),
                    nonce = oidcNonceAsHexString,
                    state = oidcNonceAsHexString
                )
                call.respond(
                    HashesSubmissionResponse(
                        idpChoices = listOf(idpRedirect.toString()),
                        salt = byteArrayToHexString(salt),
                        seed = byteArrayToHexString(seed)
                    )
                )
            }
            is Invalid -> throw input.error
        }
    }
}

