package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

sealed class Either<out T> {
    data class Error(val message: String?, val e: Exception) : Either<Nothing>()
    data class Success<T>(val value: T) : Either<T>()
}

inline fun <T> attempt(vararg args: Either<Any>, body: (args: List<Any>) -> Either<T>): Either<T> {
    args.filterIsInstance<Either.Error>().forEach { throw it.e }
    return try {
        body(args.filterIsInstance<Either.Success<Any>>().map { it.value }) as Either.Success
    } catch (e: Exception) {
        Either.Error(e.message, e)
    }
}