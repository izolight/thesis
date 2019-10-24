package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Signature

interface ISignaturesHoldingService : IExpireableCache<List<Signature>>

class SignaturesHoldingServiceDefaultImpl : ISignaturesHoldingService {
    private val expirableCache = ExpireableCacheDefaultImpl<List<Signature>>()

    override fun set(key: String, value: List<Signature>) {
        return this.set(key, value)
    }

    override fun get(key: String): List<Signature>? {
        return this.get(key)
    }

    override fun exists(key: String): Boolean {
        return this.exists(key)
    }
}
