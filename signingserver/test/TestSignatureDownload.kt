package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.URLs
import io.ktor.http.*
import io.ktor.locations.*
import io.ktor.server.testing.*
import io.ktor.util.*
import org.junit.Test
import org.koin.test.KoinTest
import kotlin.test.assertEquals

@KtorExperimentalLocationsAPI
class TestSignatureDownload : KoinTest {
    @KtorExperimentalAPI
    @Test
    fun testNonexistentSignature() {
        withTestApplication({ module() }) {
            handleRequest(HttpMethod.Get, URLs.SIGNATURES + "doesnotexist").apply {
                assertEquals(HttpStatusCode.NotFound, response.status())
            }
        }
    }
}
