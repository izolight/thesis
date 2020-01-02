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
    fun testSigningWorkflow() {


        withTestApplication({ module() }) {
            val json = Json(JsonConfiguration.Stable)
            val signatureRequest = with(handleRequest(HttpMethod.Post, URLs.SUBMIT_HASHES) {
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

                return@with SignatureRequest(
                    id_token = location.getFragments()["id_token"]
                        ?: throw IllegalArgumentException("No id_token"),
                    salt = responseBody.salt,
                    seed = responseBody.seed,
                    hashes = TESTHASHES
                )
            }

            val signatureUrl = with(handleRequest(HttpMethod.Post, URLs.SIGN) {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    json.stringify(
                        SignatureRequest.serializer(),
                        signatureRequest
                    )
                )
            }) {
                assertEquals(HttpStatusCode.OK, response.status(), response.content)
                assertNotNull(response.content)
                return@with Url(json.parse(SignatureResponse.serializer(), response.content!!).signature)
            }

            with(handleRequest(HttpMethod.Get, signatureUrl.encodedPath) {
                addHeader(HttpHeaders.Accept, ContentType.Application.OctetStream.toString())
            }) {
                assertEquals(HttpStatusCode.OK, response.status(), response.content)
                assertNotNull(response.content)
            }
        }

    }
}

