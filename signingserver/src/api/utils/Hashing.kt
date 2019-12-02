package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

import java.io.ByteArrayOutputStream
import java.security.MessageDigest

fun sha256(input: ByteArray): ByteArray = MessageDigest.getInstance("SHA-256").digest(input)

fun calculateSalt(hmacKey: ByteArray, sortedHashes: List<String>): ByteArray = hmacSha256(
    hmacKey,
    hexStringToByteArray(
        sortedHashes.joinToString("")
    )
)

fun calculateOidcNonce(maskedHashes: ByteArray): String = byteArrayToHexString(sha256(maskedHashes))

fun maskHashes(sortedHashes: List<String>, salt: ByteArray): List<ByteArray> =
    sortedHashes.map { h -> hmacSha256(salt, hexStringToByteArray(h)) }

fun List<ByteArray>.concatenate(): ByteArray = ByteArrayOutputStream(
    this.sumBy { it.size }
).also {
    this.map { array -> it.write(array) }
}.toByteArray()

