package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

import kotlinx.serialization.Serializable

@Serializable
data class SigningRequest(
    val id_token: String,
    val hashes: List<String>,
    val seed: Long,
    val salt: String
) : Validatable<SigningRequest> {
    override fun validate(): Validated<SigningRequest> {
        return when (val hashesValidationResult = SubmittedHashes(this.hashes).validate()) {
            is Valid -> Valid(
                SigningRequest(
                    this.id_token,
                    hashesValidationResult.value.hashes,
                    this.seed,
                    this.salt
                )
            )
            is Invalid -> Invalid(hashesValidationResult.error)
        }
    }
}