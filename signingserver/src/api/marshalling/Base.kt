package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

class Validated<T>(private val o: T) {
    fun get(): T {
        return o
    }
}