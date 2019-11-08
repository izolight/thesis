package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.Valid

interface IHashesCachingService : IExpireableCache<Valid<SubmittedHashes>>

class ExpirableHashesCachingServiceImpl() : IHashesCachingService {
    val expirableCache = ExpireableCacheDefaultImpl<Valid<SubmittedHashes>>()

    override fun set(key: String, value: Valid<SubmittedHashes>) {
        this.expirableCache.set(key, value)
    }

    override fun get(key: String): Valid<SubmittedHashes>? {
        return this.expirableCache.get(key)
    }

    override fun exists(key: String): Boolean {
        return this.expirableCache.exists(key)
    }
}
