package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

interface ISignaturesHoldingService : IExpireableCache<String, ByteArray> {
    fun generateId(): String
}

class SignaturesHoldingServiceDefaultImpl : ExpireableCacheDefaultImpl<String, ByteArray>(), ISignaturesHoldingService {
    private val random = SecureRandom()

    override fun generateId(): String = random.nextLong().toString(16)
}
