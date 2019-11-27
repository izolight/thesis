package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock


interface IExpireableCache<T> {
    fun set(key: String, value: T)
    fun get(key: String): T?
    fun exists(key: String): Boolean
}

class ExpireableCacheDefaultImpl<T> : IExpireableCache<T> {
    companion object {
        const val CYCLE_TIME_MILLISECONDS = 60 * 1000
        const val EXPIRATION_TIME_MILLISECONDS = 15 * 60 * 1000
    }

    data class ExpirableEntry<T>(val insertionTimeMillis: Long, val value: T)

    private var lastCycleTime = System.currentTimeMillis()
    private val storage: MutableMap<String, ExpirableEntry<T>> = HashMap()
    private val cycleLock = Mutex()

    override fun set(key: String, value: T) {
        this.storage[key] = ExpirableEntry(System.currentTimeMillis(), value)
        this.cycle()
    }

    override fun get(key: String): T? {
        this.cycle()
        return this.storage[key]?.value
    }

    override fun exists(key: String): Boolean {
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
