package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.NonceResponse
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IHashesCachingService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.INonceGeneratorService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.sha256
import io.ktor.application.call
import io.ktor.http.HttpStatusCode
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
    val hashesCache by inject<IHashesCachingService>()

    post<HashesRoute> {
        val hashes = call.receive<SubmittedHashes>()
        val randomValue = nonceGenerator.getNonce()
        val nonce = sha256(hashes.hashes + randomValue)
        hashesCache.set(nonce, hashes)
        call.respond(HttpStatusCode.Created, NonceResponse(nonce))
    }
}
