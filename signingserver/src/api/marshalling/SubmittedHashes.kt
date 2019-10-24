package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

import kotlinx.serialization.Serializable

@Serializable
data class SubmittedHashes(val hashes: List<String>) {
    init {
        if (hashes.isEmpty()) {
            throw InvalidJSONException(
                "No values provided"
            )
        }
        hashes.forEach {
            if (it.length != 64) {
                throw InvalidJSONException(
                    "Value $it is not a valid SHA256 hash: length is not 64"
                )
            }

            try {
                it.toBigInteger(16)
            } catch (e: NumberFormatException) {
                throw InvalidJSONException(
                    "Value $it is not a valid SHA256 hash: not a hex number"
                )
            }
        }
    }
}


@Serializable
data class NonceResponse(val nonce: String)