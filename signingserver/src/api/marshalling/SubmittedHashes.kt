package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

import kotlinx.serialization.Serializable

@Serializable
data class SubmittedHashes(val hashes: List<String>) : Validatable<SubmittedHashes> {
    override fun validate(): Validated<SubmittedHashes> {
        return try {
            checkForEmptyValue()
            checkForInsaneAmountsOfHashes()
            checkForValidSha256Values()
            checkForDuplicates()
            Valid(SubmittedHashes(this.hashes.sorted()))
        } catch (e: InvalidJSONException) {
            Invalid(e)
        }
    }

    private fun checkForInsaneAmountsOfHashes() {
        if(hashes.size > 100_000) {
            throw InvalidJSONException("Too many values")
        }
    }

    private fun checkForEmptyValue() {
        if (hashes.isEmpty()) {
            throw InvalidJSONException("No values provided")
        }
    }

    private fun checkForValidSha256Values() {
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

    private fun checkForDuplicates() {
        val set = HashSet<String>()
        hashes.forEach {
            if (!set.add(it)) {
                throw InvalidJSONException(
                    "Value $it was submitted more than once"
                )
            }
        }
    }
}


@Serializable
data class HashesSubmissionResponse(
    val idpChoices: List<String>,
    val seed: String,
    val salt: String
)