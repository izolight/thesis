package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import org.bouncycastle.asn1.x500.X500Name
import org.bouncycastle.crypto.CryptoException
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder
import org.bouncycastle.pkcs.PKCS10CertificationRequest
import org.bouncycastle.pkcs.jcajce.JcaPKCS10CertificationRequestBuilder
import java.security.KeyPair
import java.security.KeyPairGenerator


class Constants {
    companion object {
        val CERTIFICATE_ALGORITHM = "SHA256withRSA"
        val RSA_KEY_BITS = 4096
    }
}

class SigningKeysServiceImpl : SigningKeysService {

    private val keyPairGenerator = KeyPairGenerator.getInstance(Constants.CERTIFICATE_ALGORITHM)
    private val keyCache = ExpireableCacheDefaultImpl<KeyPair>()
    private val contentSignerBuilder = JcaContentSignerBuilder(Constants.CERTIFICATE_ALGORITHM)

    init {
        keyPairGenerator.initialize(Constants.RSA_KEY_BITS)
    }

    override fun generateSigningKey(subjectInformation: SigningKeySubjectInformation): PKCS10CertificationRequest {
        val keyPair = this.keyPairGenerator.generateKeyPair()
        this.keyCache.set(subjectInformation.toString(), keyPair)
        val csrRequestBuilder = JcaPKCS10CertificationRequestBuilder(
            X500Name(subjectInformation.toDN()),
            keyPair.public
        )
        val signer = this.contentSignerBuilder.build(keyPair.private)
        return csrRequestBuilder.build(signer) ?: throw CryptoException("Unable to construct CSR")
    }
}