package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.URLs
import io.ktor.http.HttpMethod
import io.ktor.http.HttpStatusCode
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.server.testing.handleRequest
import io.ktor.server.testing.withTestApplication
import org.junit.Test
import org.koin.test.KoinTest
import kotlin.test.assertEquals

@KtorExperimentalLocationsAPI
class TestSignatureDownload : KoinTest {
    @Test
    fun testNonexistentSignature() {
        withTestApplication({ module() }) {
            handleRequest(HttpMethod.Get, URLs.SIGNATURES + "doesnotexist").apply {
                assertEquals(HttpStatusCode.NotFound, response.status())
            }
        }
    }
}
