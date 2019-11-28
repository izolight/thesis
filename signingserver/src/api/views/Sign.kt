package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Invalid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SigningRequest
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.Config
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IOIDCService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.ISecretService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.byteArrayToHexString
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import com.auth0.jwt.interfaces.DecodedJWT
import io.ktor.application.call
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.Location
import io.ktor.locations.post
import io.ktor.request.receive
import io.ktor.routing.Routing
import org.koin.ktor.ext.inject
import org.slf4j.LoggerFactory

@KtorExperimentalLocationsAPI
@Location(URLs.SIGN)
class SignRoute

@KtorExperimentalLocationsAPI
fun Routing.sign() {
    val oidcService by inject<IOIDCService>()
    val secretService by inject<ISecretService>()
    val logger = LoggerFactory.getLogger(this.javaClass)

    fun validateOidcNonce(idToken: DecodedJWT, salt: ByteArray, concatenatedHashes: ByteArray) {
        val oidcNonce = hmacSha256(salt, concatenatedHashes)
        val oidcNonceAsHexString = byteArrayToHexString(oidcNonce)
        if (idToken.getClaim("nonce").asString() != oidcNonceAsHexString) {
            throw InvalidDataException(
                "Nonce mismatch"
            )
        }
    }

    fun validateClientId(idToken: DecodedJWT) {
        if (idToken.getClaim("aud").asString() != Config.OIDC_CLIENT_ID) {
            throw InvalidDataException("Client ID (audience) mismatch")
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



    post<SignRoute> {
        when (val input = call.receive<SigningRequest>().validate()) {
            is Valid -> {
                val salt = validateSalt(input)
                val jwtIdToken = oidcService.validateIdToken(input.value.id_token)
                validateOidcNonce(
                    jwtIdToken,
                    salt,
                    input.value.hashes.joinToString("").toByteArray(Charsets.UTF_8)
                )
                validateClientId(jwtIdToken)


            }
            is Invalid -> throw input.error
        }

    }


}

