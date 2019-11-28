package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import com.auth0.jwk.JwkException
import com.auth0.jwk.UrlJwkProvider
import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import com.auth0.jwt.exceptions.JWTDecodeException
import com.auth0.jwt.interfaces.DecodedJWT
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.client.features.json.JsonFeature
import io.ktor.client.features.json.serializer.KotlinxSerializer
import io.ktor.client.features.logging.LogLevel
import io.ktor.client.features.logging.Logging
import io.ktor.client.request.get
import io.ktor.http.Parameters
import io.ktor.http.Url
import io.ktor.http.formUrlEncode
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.Deferred
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.runBlocking
import kotlinx.serialization.Serializable
import java.net.URL
import java.security.interfaces.RSAPublicKey

interface IOIDCService {
    fun getAuthorisationEndpoint(): Url
    fun getIssuer(): Url
    fun getJwkUrl(): Url
    fun constructAuthenticationRequestUrl(
        authorisationEndpoint: Url,
        clientId: String = Config.OIDC_CLIENT_ID,
        responseType: String = Config.OIDC_RESPONSE_TYPE,
        scope: List<String> = Config.OIDC_SCOPES,
        redirectUri: Url = Config.OIDC_REDIRECT_URI,
        state: String,
        nonce: String
    ): Url

    fun validateIdToken(idToken: String): DecodedJWT
}

class Config {
    companion object {
        val OIDC_CONFIGURATION_DISCOVERY_URL = Url("https://keycloak.thesis.izolight.xyz/auth/realms/master/.well-known/openid-configuration")
        const val OIDC_CLIENT_ID = "thesis"
        const val OIDC_CLIENT_SECRET = "1f164d78-ff38-4f68-9bae-8ec8dd3b1a53"
        val OIDC_REDIRECT_URI = Url("http://127.0.0.1:8080/callback")
        val OIDC_SCOPES = listOf("openid", "profile")
        const val OIDC_RESPONSE_TYPE = "id_token"
    }
}

@KtorExperimentalAPI
class OurDemoOIDCService private constructor(
    private val futureDiscoveryDocument: Deferred<OIDCDiscoveryDocument>
) : IOIDCService {

    private val discoveryDocument: OIDCDiscoveryDocument by lazy {
        runBlocking {
            futureDiscoveryDocument.await()
        }

    }

    private val jwkProvider: UrlJwkProvider by lazy {
        UrlJwkProvider(URL(this.getJwkUrl().toString()))
    }


    companion object {
        suspend operator fun invoke() = coroutineScope {
            val client = HttpClient(CIO) {
                install(Logging) {
                    level = LogLevel.INFO
                }
                install(JsonFeature) {
                    serializer = KotlinxSerializer()
                }
            }
            val discoveryDocument = async {
                client.get<OIDCDiscoveryDocument>(Config.OIDC_CONFIGURATION_DISCOVERY_URL)
            }

            OurDemoOIDCService(futureDiscoveryDocument = discoveryDocument)
        }
    }

    override fun getAuthorisationEndpoint(): Url {
        return Url(this.discoveryDocument.authorization_endpoint)
    }

    override fun getIssuer(): Url {
        return Url(this.discoveryDocument.issuer)
    }

    override fun getJwkUrl(): Url {
        return Url(this.discoveryDocument.jwks_uri)
    }


    override fun constructAuthenticationRequestUrl(
        authorisationEndpoint: Url,
        clientId: String,
        responseType: String,
        scope: List<String>,
        redirectUri: Url,
        state: String,
        nonce: String
    ): Url {
        return Url(
            "$authorisationEndpoint?${
            Parameters.build {
                append("client_id", clientId)
                append("response_type", responseType)
                append("scope", scope.joinToString(" "))
                append("redirect_uri", redirectUri.toString())
                append("state", nonce)
                append("nonce", nonce)
            }.formUrlEncode()}"
        )
    }

    override fun validateIdToken(idToken: String): DecodedJWT {
        try {
            val jwt = JWT.decode(idToken)
            try {
                val jwk = jwkProvider.get(jwt.keyId)
                val algo = when (jwk.algorithm) {
                    "RS256" -> Algorithm.RSA256(
                        jwk.publicKey as RSAPublicKey,
                        null
                    )
                    else -> throw InvalidDataException("Unsupported algorithm")
                }
                val verifier = JWT.require(algo)
                    .withIssuer(this.getIssuer().toString())
                    .withAudience(Config.OIDC_CLIENT_ID)
                    .build()

                return verifier.verify(idToken)
            } catch (e: JwkException) {
                throw InvalidDataException(
                    "JWK Error: $e"
                )
            }
        } catch (e: JWTDecodeException) {
            throw InvalidDataException(
                "Invalid JWT: $e"
            )
        }


    }

    @Serializable
    data class OIDCDiscoveryDocument(
        val issuer: String,
        val authorization_endpoint: String,
        val token_endpoint: String,
        val jwks_uri: String,
        val subject_types_supported: List<String>,
        val response_types_supported: List<String>,
        val claims_supported: List<String>,
        val grant_types_supported: List<String>,
        val response_modes_supported: List<String>,
        val userinfo_endpoint: String,
        val scopes_supported: List<String>,
        val token_endpoint_auth_methods_supported: List<String>,
        val userinfo_signing_alg_values_supported: List<String>,
        val id_token_signing_alg_values_supported: List<String>,
        val request_parameter_supported: Boolean,
        val request_uri_parameter_supported: Boolean,
        val require_request_uri_registration: Boolean,
        val claims_parameter_supported: Boolean,
        val revocation_endpoint: String,
        val backchannel_logout_supported: Boolean,
        val backchannel_logout_session_supported: Boolean,
        val frontchannel_logout_supported: Boolean,
        val frontchannel_logout_session_supported: Boolean,
        val end_session_endpoint: String
    )
}