package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.Either
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.attempt
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.compose
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import com.auth0.jwk.Jwk
import com.auth0.jwk.UrlJwkProvider
import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import com.auth0.jwt.interfaces.DecodedJWT
import com.auth0.jwt.interfaces.JWTVerifier
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.http.Parameters
import io.ktor.http.Url
import io.ktor.http.formUrlEncode
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.Deferred
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.runBlocking
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration
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

    fun marshalJwk(jwk: Jwk): String

    fun validateIdToken(idToken: String): JwtValidationResult
    data class JwtValidationResult(val idToken: DecodedJWT, val jwk: Jwk)
}

class Config {
    companion object {
        const val OIDC_IDP_NAME = "Izolight IDP"
        val OIDC_CONFIGURATION_DISCOVERY_URL =
            Url("https://keycloak.thesis.izolight.xyz/auth/realms/master/.well-known/openid-configuration")
        const val OIDC_CLIENT_ID = "thesis"
        const val OIDC_CLIENT_SECRET = "0d6079d7-18a5-4f82-a94e-8960aed5dd89"
        val OIDC_REDIRECT_URI = Url("http://localhost:8080/callback")
        val OIDC_SCOPES = listOf("openid", "profile")
        const val OIDC_RESPONSE_TYPE = "id_token"
    }
}

@KtorExperimentalAPI
class OurDemoOIDCService private constructor(
    private val futureDiscoveryDocument: Deferred<OIDCDiscoveryDocument>
) : IOIDCService {

    @Serializable
    data class JWK(
        val kid: String,
        val kty: String,
        val alg: String,
        val use: String,
        val n: String,
        val e: String,
        val x5c: List<String>,
        val x5t: String,
        @SerialName("x5t#S256")
        val x5t_S256: String
    ) {
        companion object {
            fun fromJwk(jwks: Jwk) = JWK(
                kid = jwks.id,
                kty = jwks.type,
                alg = jwks.algorithm,
                use = jwks.usage,
                n = jwks.additionalAttributes["n"]!!.toString(),
                e = jwks.additionalAttributes["e"]!!.toString(),
                x5c = jwks.certificateChain,
                x5t = jwks.certificateThumbprint,
                x5t_S256 = jwks.additionalAttributes["x5t#S256"]!!.toString()
            )
        }
    }

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
            OurDemoOIDCService(futureDiscoveryDocument = async {
                HttpClient { defaultConfig() }.use {
                    it.get<OIDCDiscoveryDocument>(Config.OIDC_CONFIGURATION_DISCOVERY_URL)
                }
            })
        }

        private val json = Json(JsonConfiguration.Stable)
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

    override fun marshalJwk(jwk: Jwk) = json.stringify(JWK.serializer(), JWK.fromJwk(jwk))

    fun getAlgorithm(jwk: Jwk) = when (val a = jwk.algorithm) {
        "RS256" -> Algorithm.RSA256(jwk.publicKey as RSAPublicKey, null)
        else -> throw InvalidDataException("Unsupported algorithm in JWK: $a")
    }

    fun buildVerifier(algorithm: Algorithm): JWTVerifier = JWT.require(algorithm)
        .withIssuer(this.getIssuer().toString())
        .withAudience(Config.OIDC_CLIENT_ID)
        .build()

    override fun validateIdToken(idToken: String): IOIDCService.JwtValidationResult {
        val result = compose {
            val jwt = attempt { JWT.decode(idToken) }
            val jwk = attempt(jwt) { jwkProvider.get(it.keyId) }
            val algo = attempt(jwk) { getAlgorithm(it) }
            val verifier = attempt(algo) { buildVerifier(it) }
            attempt(verifier, jwk) { v, j ->
                IOIDCService.JwtValidationResult(idToken = v.verify(idToken), jwk = j)
            }
        }
        when(result) {
            is Either.Success -> return result.value
            is Either.Error -> throw result.e
        }
    }

    @Serializable
    data class OIDCDiscoveryDocument(
        val issuer: String,
        val authorization_endpoint: String,
        val token_endpoint: String,
        val token_introspection_endpoint: String,
        val userinfo_endpoint: String,
        val end_session_endpoint: String,
        val jwks_uri: String,
        val check_session_iframe: String,
        val grant_types_supported: List<String>,
        val response_types_supported: List<String>,
        val subject_types_supported: List<String>,
        val id_token_signing_alg_values_supported: List<String>,
        val id_token_encryption_alg_values_supported: List<String>,
        val id_token_encryption_enc_values_supported: List<String>,
        val userinfo_signing_alg_values_supported: List<String>,
        val request_object_signing_alg_values_supported: List<String>,
        val response_modes_supported: List<String>,
        val registration_endpoint: String,
        val token_endpoint_auth_methods_supported: List<String>,
        val token_endpoint_auth_signing_alg_values_supported: List<String>,
        val claims_supported: List<String>,
        val claim_types_supported: List<String>,
        val claims_parameter_supported: Boolean,
        val scopes_supported: List<String>,
        val request_parameter_supported: Boolean,
        val request_uri_parameter_supported: Boolean,
        val code_challenge_methods_supported: List<String>,
        val tls_client_certificate_bound_access_tokens: Boolean,
        val introspection_endpoint: String
    )
}