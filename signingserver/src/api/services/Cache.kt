package ch.bfh.ti.hirtp1ganzg1.thesis.api.services


interface ICache<T, U> {
    fun set(key: T, value: U)
    fun get(key: T): U?
    fun remove(key: T)
    fun exists(key: T): Boolean
}

