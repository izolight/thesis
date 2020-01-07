package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def.ISecretService
import org.bouncycastle.crypto.digests.SHA256Digest
import org.bouncycastle.crypto.generators.HKDFBytesGenerator
import org.bouncycastle.crypto.params.HKDFParameters
import java.security.SecureRandom


class SecretServiceDefaultImpl : ISecretService {
    private val secret = ByteArray(64).also { SecureRandom().nextBytes(it) }
    override fun hkdf(salt: ByteArray, length: Int): ByteArray = ByteArray(length).also {
        HKDFBytesGenerator(SHA256Digest()).also { hkdf ->
            hkdf.init(
                HKDFParameters(this.secret, salt, ByteArray(0))
            )
            hkdf.generateBytes(it, 0, length)
        }
    }
}