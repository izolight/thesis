package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.Constants
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
import io.ktor.http.Url
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.Deferred
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.runBlocking
import java.net.URL
import java.net.URLEncoder
import java.security.interfaces.RSAPublicKey

interface IOIDCService {
    fun getAuthorisationEndpoint(): Url
    fun getIssuer(): Url
    fun getJwkUrl(): Url
    fun constructAuthenticationRequestUrl(
        authorisationEndpoint: Url,
        clientId: String = Constants.OIDC_CLIENT_ID,
        responseType: String = Constants.OIDC_RESPONSE_TYPE,
        scope: List<String> = Constants.OIDC_SCOPES,
        redirectUri: Url = Constants.OIDC_REDIRECT_URI,
        state: String,
        nonce: String
    ): Url

    fun validateIdToken(idToken: String): DecodedJWT
}

@KtorExperimentalAPI
class GoogleOIDCService private constructor(
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
            val discoveryDocumentUrl = Url("https://accounts.google.com/.well-known/openid-configuration")
            val client = HttpClient(CIO) {
                install(Logging) {
                    level = LogLevel.INFO
                }
                install(JsonFeature) {
                    serializer = KotlinxSerializer()
                }
            }
            val discoveryDocument = async {
                client.get<OIDCDiscoveryDocument>(discoveryDocumentUrl)
            }

            GoogleOIDCService(futureDiscoveryDocument = discoveryDocument)
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
        val scopeConcatenated = scope.fold("", { acc, next -> "$acc $next" })
        return Url(
            URLEncoder.encode(
                "$authorisationEndpoint?client_id=$clientId&response_type=$responseType&scope=$scopeConcatenated&redirect_uri=$redirectUri&state=$state&nonce=$nonce",
                "UTF-8"
            )
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
                    .withAudience(Constants.OIDC_CLIENT_ID)
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

    data class OIDCDiscoveryDocument(
        val issuer: String,
        val authorization_endpoint: String,
        val token_endpoint: String,
        val userinfo_endpoint: String,
        val revocation_endpoint: String,
        val jwks_uri: String
    )
}