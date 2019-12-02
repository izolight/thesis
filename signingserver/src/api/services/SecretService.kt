package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import org.bouncycastle.crypto.digests.SHA256Digest
import org.bouncycastle.crypto.generators.HKDFBytesGenerator
import org.bouncycastle.crypto.params.HKDFParameters
import java.security.SecureRandom

interface ISecretService {
//    fun getSecret(): ByteArray
    fun hkdf(salt: ByteArray, length: Int = 64): ByteArray
}


class SecretServiceDefaultImpl : ISecretService {
    private val secret = ByteArray(32).also { SecureRandom().nextBytes(it) }
//    override fun getSecret(): ByteArray = this.secret.copyOf()
    override fun hkdf(salt: ByteArray, length: Int): ByteArray = ByteArray(length).also {
        HKDFBytesGenerator(SHA256Digest()).also { hkdf ->
            hkdf.init(
                HKDFParameters(this.secret, salt, ByteArray(0))
            )
            hkdf.generateBytes(it, 0, length)
        }
    }
}