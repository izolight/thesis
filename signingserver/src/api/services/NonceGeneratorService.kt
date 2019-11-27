package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

interface INonceGeneratorService {
    fun getNonce(): ByteArray
}

class NonceGeneratorServiceDefaultImpl : INonceGeneratorService {
    private val secureRandom = SecureRandom()

    override fun getNonce(): ByteArray {
        return ByteArray(10).also {
            secureRandom.nextBytes(it)
        }
    }
}