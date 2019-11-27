package ch.bfh.ti.hirtp1ganzg1.thesis

import at.favre.lib.bytes.Bytes
import io.ktor.http.HttpMethod
import io.ktor.http.HttpStatusCode
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.server.testing.handleRequest
import io.ktor.server.testing.withTestApplication
import org.koin.test.KoinTest
import kotlin.test.Test
import kotlin.test.assertEquals

@KtorExperimentalLocationsAPI
class ApplicationTest : KoinTest {
    @Test
    fun testByteSerialisation() {
        val someString = "hello I would like to serialise plix"
        val byteArray = someString.toByteArray()
        val hexString = Bytes.wrap(byteArray).encodeHex()
        val byteArray2 = Bytes.parseHex(hexString).array()
        byteArray.forEachIndexed { index, byte ->  assertEquals(byte, byteArray2[index]) }
    }

    @Test
    fun testRoot() {
        withTestApplication({ module() }) {
            handleRequest(HttpMethod.Get, "/").apply {
                assertEquals(HttpStatusCode.OK, response.status())
                assertEquals("lol generics", response.content)
            }
        }
    }

}

