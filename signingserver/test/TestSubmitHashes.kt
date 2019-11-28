package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.response.HttpResponse
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
        @Serializable
        data class TestSubmitHashesPostBody(val hashes: List<String>)

        @Serializable
        data class ExpectedNonceResponse(val idpChoices: List<String>, val seed: String, val salt: String)

        withTestApplication({ module() }) {
            val json = Json(JsonConfiguration.Stable)
            with(handleRequest(HttpMethod.Post, "/api/v1/hashes") {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    json.stringify(
                        TestSubmitHashesPostBody.serializer(),
                        TestSubmitHashesPostBody(
                            listOf(
                                "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308",
                                "8f33ddf44093ee0cc72c7123f878a8926feab6cedf885e148d45ae30213cd443"
                            )
                        )
                    )
                )
            }) {
                assertEquals(HttpStatusCode.OK, response.status(), response.content)
                val responseText = response.content.toString()
                assertTrue("nonce" in responseText, responseText)
                val response = json.parse(ExpectedNonceResponse.serializer(), responseText)
                assertNotNull(response)
                assertFalse(response.idpChoices.isEmpty())
                response.idpChoices.forEach { s -> Url(s) }

                runBlocking {
                    val client = HttpClient() { defaultConfig() }
                    val idpResponse = client.get<HttpResponse>(response.idpChoices[0])
                    assertEquals(HttpStatusCode.OK, idpResponse.status)
//                    TODO("submit form with csrf and challenge stuff")
//                    val loginResponse = client.post<HttpResponse>(
//                        "https://${idpResponse.call.request.url.hostWithPort}/login",
//                        MultiPartFormDataContent(formData {
//
//                        })
//                    )
                }
            }

            with(handleRequest(HttpMethod.Post, "/api/v1/hashes") {
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

            with(handleRequest(HttpMethod.Post, "/api/v1/hashes") {
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
