package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

import at.favre.lib.bytes.Bytes
import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec

class Constants {
    companion object {
        const val HMAC_SHA256 = "HmacSHA256"
    }
}

fun hmacSha256(key: ByteArray, contents: ByteArray): ByteArray {
    val keySpec = SecretKeySpec(key, Constants.HMAC_SHA256)
    val mac = Mac.getInstance(Constants.HMAC_SHA256)
    mac.init(keySpec)
    return mac.doFinal(contents)
}

fun byteArrayToHexString(a: ByteArray): String = Bytes.wrap(a).encodeHex()

fun hexStringToByteArray(s: String): ByteArray = Bytes.parseHex(s).array()
