package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

interface INonceGeneratorService {
    fun getNonce(): Long
}

class NonceGeneratorServiceDefaultImpl : INonceGeneratorService {
    private val secureRandom = SecureRandom()

    override fun getNonce(): Long {
        return secureRandom.nextLong()
    }
}