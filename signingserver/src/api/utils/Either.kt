package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

sealed class Either<out T> {
    data class Error(val message: String?, val e: Exception) : Either<Nothing>()
    data class Success<T>(val value: T) : Either<T>()
}

inline fun <T> compose(body: () -> Either<T>): Either<T> {
    return try {
        body()
    } catch (e: Exception) {
        Either.Error(e.message, e)
    }
}

inline fun <T> attempt(body: () -> T): Either<T> {
    return try {
        Either.Success(body())
    } catch (e: Exception) {
        Either.Error(e.message, e)
    }
}

inline fun <T, U> attempt(a: Either<T>, body: (a: T) -> U): Either<U> {
    when (a) {
        is Either.Error -> throw a.e
        is Either.Success ->
            return try {
                Either.Success(body(a.value))
            } catch (e: Exception) {
                Either.Error(e.message, e)
            }
    }
}


inline fun <T, U, V> attempt(a: Either<T>, b: Either<U>, body: (a: T, b: U) -> V): Either<V> {
    when (a) {
        is Either.Error -> throw a.e
        is Either.Success -> when (b) {
            is Either.Error -> throw b.e
            is Either.Success -> return try {
                Either.Success(body(a.value, b.value))
            } catch (e: Exception) {
                Either.Error(e.message, e)
            }
        }
    }
}

inline fun <T, U, V, W> attempt(
    a: Either<T>,
    b: Either<U>,
    c: Either<V>,
    body: (a: T, b: U, c: V) -> W
): Either<W> {
    when (a) {
        is Either.Error -> throw a.e
        is Either.Success -> when (b) {
            is Either.Error -> throw b.e
            is Either.Success -> when (c) {
                is Either.Error -> throw c.e
                is Either.Success -> return try {
                    Either.Success(body(a.value, b.value, c.value))
                } catch (e: Exception) {
                    Either.Error(e.message, e)
                }
            }
        }
    }
}
