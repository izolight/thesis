package ch.bfh.ti.hirtp1ganzg1.thesis

import io.ktor.http.Url
import kotlinx.serialization.Serializable

val TESTUSERNAME = "testuser2"
val TESTPASSWORD = "test1234"
val TESTHASHES = listOf(
    "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308",
    "8f33ddf44093ee0cc72c7123f878a8926feab6cedf885e148d45ae30213cd443"
)

@Serializable
data class TestSubmitHashesPostBody(val hashes: List<String>)

@Serializable
data class ExpectedNonceResponse(val providers: Map<String, String>, val seed: String, val salt: String)

@Serializable
data class SignatureRequest(
    val id_token: String,
    val seed: String,
    val salt: String,
    val hashes: List<String>
)

@Serializable
data class SignatureResponse(
    val signature: String
)

fun Url.getFragments(): HashMap<String, String> = HashMap<String, String>().also {
    fragment.splitToSequence("&").forEach { parameter ->
        parameter.split("=").also { keyValuePair ->
            it[keyValuePair[0]] = keyValuePair[1]
        }
    }
}
