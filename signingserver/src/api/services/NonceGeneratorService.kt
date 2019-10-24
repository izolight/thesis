package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

object NonceGeneratorService {
    private val secureRandom = SecureRandom()

    fun getNonce(): String {
        return secureRandom.nextLong().toString(16)
    }
}