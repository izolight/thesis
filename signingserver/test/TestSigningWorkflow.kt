package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl.Config
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.URLs
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.request.*
import io.ktor.client.request.forms.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.serialization.*
import io.ktor.server.testing.*
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.runBlocking
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.encodeToString
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
            val signatureRequest = with(handleRequest(HttpMethod.Post, URLs.SUBMIT_HASHES) {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    DefaultJson.encodeToString(
                        TestSubmitHashesPostBody(TESTHASHES)
                    )
                )
            }) {
                assertEquals(HttpStatusCode.OK, response.status(), response.content)
                val responseText = response.content.toString()
                assertTrue("nonce" in responseText, responseText)
                val responseBody = DefaultJson.decodeFromString<ExpectedNonceResponse>(responseText)
                assertNotNull(responseBody)
                assertFalse(responseBody.providers.isEmpty())
                assertTrue(responseBody.providers.containsKey(Config.OIDC_IDP_NAME))
                val idpUrl = responseBody.providers[Config.OIDC_IDP_NAME]
                assertNotNull(idpUrl)
                responseBody.providers.forEach { entry -> Url(entry.value) }

                val location = runBlocking {
                    val client = HttpClient(CIO) { defaultConfig().also { followRedirects = false } }
                    val initialIdpResponse = client.get<HttpStatement>(idpUrl).execute()
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
                    DefaultJson.encodeToString(
                        SignatureRequest.serializer(),
                        signatureRequest
                    )
                )
            }) {
                assertEquals(HttpStatusCode.OK, response.status(), response.content)
                assertNotNull(response.content)
                return@with Url(DefaultJson.decodeFromString<SignatureResponse>(response.content!!).signature)
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

