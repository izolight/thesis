package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def.ICache

open class ICacheDefaultImpl<T, U> : ICache<T, U> {
    private val storage = HashMap<T, U>()

    override fun set(key: T, value: U) = this.storage.set(key, value)

    override fun get(key: T): U? = this.storage[key]

    override fun remove(key: T) {
        this.storage.remove(key)
    }

    override fun exists(key: T): Boolean = this.storage.containsKey(key)
}