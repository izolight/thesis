package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.IHashesCachingService
import com.beust.klaxon.Klaxon
import io.ktor.http.ContentType
import io.ktor.http.HttpHeaders
import io.ktor.http.HttpMethod
import io.ktor.http.HttpStatusCode
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.server.testing.handleRequest
import io.ktor.server.testing.setBody
import io.ktor.server.testing.withTestApplication
import org.koin.test.KoinTest
import org.koin.test.inject
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

@KtorExperimentalLocationsAPI
class ApplicationTest : KoinTest {
    @Test
    fun testRoot() {
        withTestApplication({ module() }) {
            handleRequest(HttpMethod.Get, "/").apply {
                assertEquals(HttpStatusCode.OK, response.status())
                assertEquals("lol generics", response.content)
            }
        }
    }

    @Test
    fun testSubmitHashes() {
        data class TestSubmitHashesPostBody(val hashes: List<String>)
        data class ExpectedNonceResponse(val nonce: String)

        val klaxon = Klaxon()
        withTestApplication({ module() }) {
            val hashesCache by inject<IHashesCachingService>()
            with(handleRequest(HttpMethod.Post, "/api/v1/hashes") {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    klaxon.toJsonString(
                        TestSubmitHashesPostBody(
                            listOf(
                                "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308",
                                "8f33ddf44093ee0cc72c7123f878a8926feab6cedf885e148d45ae30213cd443"
                            )
                        )
                    )
                )
            }) {
                assertEquals(HttpStatusCode.Created, response.status())
                val responseText = response.content.toString()
                assertTrue("nonce" in responseText, responseText)
                val response = klaxon.parse<ExpectedNonceResponse>(responseText)
                assertNotNull(response)
                assertTrue(response.nonce.length == 64)
                assertTrue(hashesCache.exists(response.nonce))
            }

            with(handleRequest(HttpMethod.Post, "/api/v1/hashes") {
                addHeader(HttpHeaders.ContentType, ContentType.Application.Json.toString())
                addHeader(HttpHeaders.Accept, ContentType.Application.Json.toString())

                setBody(
                    klaxon.toJsonString(
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
                    klaxon.toJsonString(
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

