package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock


interface IExpireableCache<T, U> {
    fun set(key: T, value: U)
    fun get(key: T): U?
    fun exists(key: T): Boolean
}

class ExpireableCacheDefaultImpl<T, U> : IExpireableCache<T, U> {
    companion object {
        const val CYCLE_TIME_MILLISECONDS = 60 * 1000
        const val EXPIRATION_TIME_MILLISECONDS = 15 * 60 * 1000
    }

    data class ExpirableEntry<U>(val insertionTimeMillis: Long, val value: U)

    private var lastCycleTime = System.currentTimeMillis()
    private val storage: MutableMap<T, ExpirableEntry<U>> = HashMap()
    private val cycleLock = Mutex()

    override fun set(key: T, value: U) {
        this.storage[key] = ExpirableEntry(System.currentTimeMillis(), value)
        this.cycle()
    }

    override fun get(key: T): U? {
        this.cycle()
        return this.storage[key]?.value
    }

    override fun exists(key: T): Boolean {
        return this.storage.contains(key)
    }

    private fun isTimeToCycle(): Boolean {
        return System.currentTimeMillis() > this.lastCycleTime + CYCLE_TIME_MILLISECONDS
    }

    private fun cycle() {
        if (isTimeToCycle() && !cycleLock.isLocked) {
            GlobalScope.launch {
                cycleLock.withLock {
                    val now = System.currentTimeMillis()
                    storage.forEach { (key, value) ->
                        if (value.insertionTimeMillis + EXPIRATION_TIME_MILLISECONDS > now) {
                            storage.remove(key)
                        }
                    }
                    lastCycleTime = System.currentTimeMillis()
                }
            }
        }
    }
}
