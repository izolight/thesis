package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

import kotlinx.serialization.Serializable

@Serializable
data class SigningRequest(
    val id_token: String,
    val seed: String,
    val salt: String,
    val hashes: List<String>
) : Validatable<SigningRequest> {
    override fun validate(): Validated<SigningRequest> {
        if (seed.length != 64) return Invalid(InvalidJSONException("Invalid seed length"))
        if (salt.length != 64) return Invalid(InvalidJSONException("Invalid salt length"))

        return when (val hashesValidationResult = SubmittedHashes(this.hashes).validate()) {
            is Valid -> Valid(
                SigningRequest(
                    this.id_token,
                    this.seed,
                    this.salt,
                    hashesValidationResult.value.hashes
                )
            )
            is Invalid -> Invalid(hashesValidationResult.error)
        }
    }
}

@Serializable
data class SigningResponse(
    val signature: String
)