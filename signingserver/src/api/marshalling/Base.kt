package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling


sealed class Validated<T>
class Valid<T>(val value: T) : Validated<T>()
class Invalid<T>(val error: Exception) : Validated<T>()

interface Validatable<T> {
    fun validate(): Validated<T>
}
