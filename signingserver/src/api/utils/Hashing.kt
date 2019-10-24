package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

import java.security.MessageDigest

fun sha256(values: List<String>): String {
    val bytes = values.joinToString().toByteArray()
    val digester = MessageDigest.getInstance("SHA-256")
    val digest = digester.digest(bytes)
    return digest.fold("", { acc, it -> acc + "%02x".format(it) })
}