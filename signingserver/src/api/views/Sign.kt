package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SigningRequest
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IOIDCService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.ISecretService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hexStringToByteArray
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import com.auth0.jwt.interfaces.DecodedJWT
import io.ktor.application.call
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.Location
import io.ktor.locations.post
import io.ktor.request.receive
import io.ktor.routing.Routing
import org.koin.ktor.ext.inject

@KtorExperimentalLocationsAPI
@Location("/api/v1/sign")
class SignRoute()

@KtorExperimentalLocationsAPI
fun Routing.sign() {
    val oidcService by inject<IOIDCService>()
    val secretService by inject<ISecretService>()

    fun validateOidcNonceAndState(idToken: DecodedJWT, salt: ByteArray, concatenatedHashes: ByteArray) {
        val oidcNonce = hmacSha256(salt, concatenatedHashes).toString()
        if (idToken.getClaim("nonce").asString() != oidcNonce) {
            throw InvalidDataException(
                "Nonce mismatch"
            )
        } else {
            if (idToken.getClaim("state").asString() != oidcNonce) {
                throw InvalidDataException(
                    "State mismatch"
                )
            }
        }
    }

    fun validateSalt(signingRequest: Valid<SigningRequest>) {
        val hmacKey = hexStringToByteArray(signingRequest.value.seed)
        val concatenatedHashes = signingRequest.value.hashes.fold(
            "",
            { acc, next ->
                acc + next
            }
        ).toByteArray()
        val salt = hmacSha256(hmacKey, concatenatedHashes)

        if (salt.toString() != signingRequest.value.salt) {
            throw InvalidDataException(
                "Salt mismatch"
            )
        } else {
            val jwtIdToken = oidcService.validateIdToken(signingRequest.value.id_token)
            validateOidcNonceAndState(
                jwtIdToken,
                salt,
                concatenatedHashes
            )
        }
    }



    post<SignRoute> {
        when (val input = call.receive<SigningRequest>().validate()) {
            is Valid -> {
                validateSalt(input)


            }
            is Invalid -> throw input.error
        }

    }


}

