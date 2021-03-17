package ch.bfh.ti.hirtp1ganzg1.thesis

import at.favre.lib.bytes.Bytes
import io.ktor.http.*
import io.ktor.locations.*
import io.ktor.server.testing.*
import io.ktor.util.*
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

    @KtorExperimentalAPI
    @Test
    fun testRoot() {
        withTestApplication({ module() }) {
            handleRequest(HttpMethod.Get, "/").apply {
                assertEquals(HttpStatusCode.Found, response.status())
            }
        }
    }

}

