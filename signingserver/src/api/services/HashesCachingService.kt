package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.SubmittedHashes

interface IHashesCachingService : IExpireableCache<SubmittedHashes>

class ExpirableHashesCachingServiceImpl() : IHashesCachingService {
    val expirableCache = ExpireableCacheDefaultImpl<SubmittedHashes>()

    override fun set(key: String, value: SubmittedHashes) {
        this.expirableCache.set(key, value)
    }

    override fun get(key: String): SubmittedHashes? {
        return this.expirableCache.get(key)
    }

    override fun exists(key: String): Boolean {
        return this.expirableCache.exists(key)
    }
}
