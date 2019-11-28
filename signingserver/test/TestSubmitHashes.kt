package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.Config
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.URLs
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.client.request.forms.submitForm
import io.ktor.client.request.get
import io.ktor.client.request.header
import io.ktor.client.response.HttpResponse
import io.ktor.client.response.readText
import io.ktor.http.*
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.server.testing.handleRequest
import io.ktor.server.testing.setBody
import io.ktor.server.testing.withTestApplication
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.runBlocking
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration
import org.jsoup.Jsoup
import org.jsoup.nodes.FormElement
import org.junit.Test
import org.koin.test.KoinTest
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

@KtorExperimentalLocationsAPI
class TestSubmitHashes : KoinTest {
    @KtorExperimentalAPI
    @Test
    fun testSubmitHashes() {

        val TESTUSERNAME = "testuser"
        val TESTPASSWORD = "test1234"
        val TESTHASHES = listOf(
            "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308",
            "8f33ddf44093ee0cc72c7123f878a8926feab6cedf885e148d45ae30213cd443"
        )

        @Serializable
        data class TestSubmitHashesPostBody(val hashes: List<String>)

        @Serializable
        data class ExpectedNonceResponse(val providers: Map<String, String>, val seed: String, val salt: String)

        withTestApplication({ module() }) {
            val json = Json(JsonConfiguration.Stable)
            with(handleRequest(HttpMethod.Post, URLs.SUBMIT_HASHES) {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    json.stringify(
                        TestSubmitHashesPostBody.serializer(),
                        TestSubmitHashesPostBody(TESTHASHES)
                    )
                )
            }) {
                assertEquals(HttpStatusCode.OK, response.status(), response.content)
                val responseText = response.content.toString()
                assertTrue("nonce" in responseText, responseText)
                val responseBody = json.parse(ExpectedNonceResponse.serializer(), responseText)
                assertNotNull(responseBody)
                assertFalse(responseBody.providers.isEmpty())
                assertTrue(responseBody.providers.containsKey(Config.OIDC_IDP_NAME))
                val idpUrl = responseBody.providers[Config.OIDC_IDP_NAME]
                assertNotNull(idpUrl)
                responseBody.providers.forEach { entry -> Url(entry.value) }

                val location = runBlocking {
                    val client = HttpClient(CIO) { defaultConfig().also { followRedirects = false } }
                    val initialIdpResponse = client.get<HttpResponse>(idpUrl)
                    assertEquals(initialIdpResponse.status, HttpStatusCode.OK)

                    val htmlLoginForm =
                        (Jsoup.parse(initialIdpResponse.readText(Charsets.UTF_8)).getElementById("kc-form-login")!! as FormElement)
                    val formTargetUrl = Url(htmlLoginForm.attributes().get("action"))
                    val idpToSigningServiceCallback = client.submitForm<HttpResponse>(url = formTargetUrl.toString(),
                        formParameters = Parameters.build {
                            append("username", TESTUSERNAME)
                            append("password", TESTPASSWORD)
                            append("credentialId", "")
                        }) {
                        method = HttpMethod.Post
                        header("Cookie", initialIdpResponse.headers["Set-Cookie"])
                    }
                    assertEquals(idpToSigningServiceCallback.status, HttpStatusCode.Found)
                    assertTrue(idpToSigningServiceCallback.headers.contains("Location"))

                    return@runBlocking Url(idpToSigningServiceCallback.headers["Location"]!!)
                }

                @Serializable
                data class SignatureRequest(
                    val id_token: String,
                    val seed: String,
                    val salt: String,
                    val hashes: List<String>
                )

                with(handleRequest(HttpMethod.Post, URLs.SIGN) {
                    addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                    addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                    setBody(
                        json.stringify(
                            SignatureRequest.serializer(),
                            SignatureRequest(
                                id_token = location.getFragments()["id_token"]
                                    ?: throw IllegalArgumentException("No id_token"),
                                salt = responseBody.salt,
                                seed = responseBody.seed,
                                hashes = TESTHASHES
                            )
                        )
                    )
                }) {
                    TODO("split up unit tests - avoid nested with(handleRequest())")
                    assertEquals(HttpStatusCode.OK, response.status(), response.content)
                    println()
                }

                with(handleRequest(HttpMethod.Post, URLs.SIGN) {
                    addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                    addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                    setBody(
                        json.stringify(
                            SignatureRequest.serializer(),
                            SignatureRequest(
                                id_token = "${location.getFragments()["id_token"]
                                    ?: throw IllegalArgumentException("No id_token")}umad",
                                salt = responseBody.salt,
                                seed = responseBody.seed,
                                hashes = TESTHASHES
                            )
                        )
                    )
                }) {
                    assertEquals(HttpStatusCode.BadRequest, response.status(), response.content)
                }

                with(handleRequest(HttpMethod.Post, URLs.SIGN) {
                    addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                    addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                    setBody(
                        json.stringify(
                            SignatureRequest.serializer(),
                            SignatureRequest(
                                id_token = location.getFragments()["id_token"]
                                    ?: throw IllegalArgumentException("No id_token"),
                                salt = responseBody.salt + "a",
                                seed = responseBody.seed,
                                hashes = TESTHASHES
                            )
                        )
                    )
                }) {
                    assertEquals(HttpStatusCode.BadRequest, response.status(), response.content)
                }

                with(handleRequest(HttpMethod.Post, URLs.SIGN) {
                    addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                    addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                    setBody(
                        json.stringify(
                            SignatureRequest.serializer(),
                            SignatureRequest(
                                id_token = location.getFragments()["id_token"]
                                    ?: throw IllegalArgumentException("No id_token"),
                                salt = responseBody.salt + "a",
                                seed = responseBody.seed.replace(responseBody.seed[0], responseBody.seed[0] + 1),
                                hashes = TESTHASHES
                            )
                        )
                    )
                }) {
                    assertEquals(HttpStatusCode.BadRequest, response.status(), response.content)
                }


            }


            with(handleRequest(HttpMethod.Post, URLs.SUBMIT_HASHES) {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    json.stringify(
                        TestSubmitHashesPostBody.serializer(),
                        TestSubmitHashesPostBody(
                            listOf(
                                "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308",
                                "8f33ddf43ee0cc72c7123f878a8926feab6cedf885e148d45ae30213cd443"
                            )
                        )
                    )
                )
            }) {
                assertEquals(
                    HttpStatusCode.BadRequest,
                    response.status(),
                    "Status: ${response.status().toString()}, body: ${response.content}"
                )
                val responseText = response.content.toString()
                assertTrue("not a valid" in responseText, responseText)
            }

            with(handleRequest(HttpMethod.Post, URLs.SUBMIT_HASHES) {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    json.stringify(
                        TestSubmitHashesPostBody.serializer(),
                        TestSubmitHashesPostBody(
                            listOf(
                            )
                        )
                    )
                )
            }) {
                assertEquals(HttpStatusCode.BadRequest, response.status())
                val responseText = response.content.toString()
                assertTrue("No values" in responseText, responseText)
            }
        }
    }
}

fun Url.getFragments(): HashMap<String, String> {
    return HashMap<String, String>().also {
        fragment.splitToSequence("&").forEach { parameter ->
            parameter.split("=").also { keyValuePair ->
                it[keyValuePair[0]] = keyValuePair[1]
            }
        }
    }
}
