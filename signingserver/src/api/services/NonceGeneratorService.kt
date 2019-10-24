package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

interface INonceGeneratorService {
    fun getNonce(): String
}

class NonceGeneratorServiceDefaultImpl : INonceGeneratorService {
    private val secureRandom = SecureRandom()

    override fun getNonce(): String {
        return secureRandom.nextLong().toString(16)
    }
}