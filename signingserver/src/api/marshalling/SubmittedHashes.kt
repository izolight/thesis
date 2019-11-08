package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

import kotlinx.serialization.Serializable

@Serializable
data class SubmittedHashes(val hashes: List<String>) : Validatable<SubmittedHashes> {
    override fun validate(): Validated<SubmittedHashes> {
        if (hashes.isEmpty()) {
            return Invalid(InvalidJSONException("No values provided"))
        }
        hashes.forEach {
            if (it.length != 64) {
                return Invalid(
                    InvalidJSONException(
                        "Value $it is not a valid SHA256 hash: length is not 64"
                    )
                )
            }

            try {
                it.toBigInteger(16)
            } catch (e: NumberFormatException) {
                return Invalid(
                    InvalidJSONException(
                        "Value $it is not a valid SHA256 hash: not a hex number"
                    )
                )
            }
        }

        return Valid(this)
    }
}


@Serializable
data class NonceResponse(val nonce: String)