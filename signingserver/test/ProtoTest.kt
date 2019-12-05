package ch.bfh.ti.hirtp1ganzg1.thesis

import Signature
import com.google.protobuf.ByteString
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue


//@Test
class ProtoTest {
    @Test
    fun serialiseThenDeserialiseTest() {
        val msg = Signature.SignatureData.newBuilder()
            .setHashAlgorithm(Signature.HashAlgorithm.SHA2_256)
            .setIdToken(ByteString.copyFrom("fubar", Charsets.UTF_8))
            .setMacAlgorithm(Signature.MACAlgorithm.HMAC_SHA2_256)
            .build()

        val serialisedMsg = msg.toByteArray()
        assertTrue(serialisedMsg.isNotEmpty())

        val msg2 = Signature.SignatureData.parseFrom(serialisedMsg)
        assertEquals(msg2.idToken, ByteString.copyFrom("fubar", Charsets.UTF_8))
        assertEquals(msg2.hashAlgorithm, Signature.HashAlgorithm.SHA2_256)
        assertEquals(msg2.macAlgorithm, Signature.MACAlgorithm.HMAC_SHA2_256)
    }
}