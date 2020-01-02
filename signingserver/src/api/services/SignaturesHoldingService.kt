package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom
import kotlin.math.absoluteValue

interface ISignaturesHoldingService : ICache<String, ByteArray> {
    fun generateId(): String
}

class SignaturesHoldingServiceDefaultImpl : ICacheDefaultImpl<String, ByteArray>(), ISignaturesHoldingService {
    private val random = SecureRandom()

    override fun generateId(): String = random.nextLong().absoluteValue.toString(16)
}
