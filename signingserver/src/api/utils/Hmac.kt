package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec

class Constants {
    companion object {
        val HMAC_SHA256 = "HmacSHA256"
    }
}

fun hmacSha256(key: ByteArray, contents: ByteArray): ByteArray {
    val keySpec = SecretKeySpec(key, Constants.HMAC_SHA256)
    val mac = Mac.getInstance(Constants.HMAC_SHA256)
    mac.init(keySpec)
    return mac.doFinal(contents)
}