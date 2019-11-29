package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import Signature
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers
import org.bouncycastle.asn1.x500.X500Name
import org.bouncycastle.asn1.x509.*
import org.bouncycastle.cms.CMSSignedData
import org.bouncycastle.crypto.CryptoException
import org.bouncycastle.jce.X509KeyUsage
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder
import org.bouncycastle.pkcs.PKCS10CertificationRequest
import org.bouncycastle.pkcs.jcajce.JcaPKCS10CertificationRequestBuilder
import java.security.KeyPair
import java.security.KeyPairGenerator


class Constants {
    companion object {
        const val KEY_ALGORITHM = "RSA"
        const val SIGNATURE_ALGORITHM = "SHA256withRSA"
        const val RSA_KEY_BITS = 4096
    }
}

class SigningKeysServiceImpl : ISigningKeysService {
    private val keyPairGenerator = KeyPairGenerator.getInstance(Constants.KEY_ALGORITHM)
    private val keyCache = ExpireableCacheDefaultImpl<KeyPair>()
    private val contentSignerBuilder = JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)

    init {
        keyPairGenerator.initialize(Constants.RSA_KEY_BITS)
    }

    override fun generateSigningKey(subjectInformation: SigningKeySubjectInformation): PKCS10CertificationRequest {
        val keyPair = this.keyPairGenerator.generateKeyPair()
        this.keyCache.set(subjectInformation.toString(), keyPair)

        return JcaPKCS10CertificationRequestBuilder(
            X500Name(subjectInformation.toDN()),
            keyPair.public
        ).setLeaveOffEmptyAttributes(
            true
        ).addAttribute(
            PKCSObjectIdentifiers.pkcs_9_at_extensionRequest,
            ExtensionsGenerator().also {
                it.addExtension(
                    Extension.basicConstraints,
                    true,
                    BasicConstraints(false)
                )
                it.addExtension(
                    Extension.keyUsage,
                    true,
                    X509KeyUsage(X509KeyUsage.digitalSignature).toASN1Primitive()
                )
                it.addExtension(
                    Extension.subjectAlternativeName,
                    false,
                    GeneralNames(GeneralName(GeneralName.rfc822Name, subjectInformation.email))
                )
            }.generate()
        ).build(this.contentSignerBuilder.build(keyPair.private)) ?: throw CryptoException("Unable to construct CSR")
    }

    override fun signToPkcs7(signatureData: Signature.SignatureData): CMSSignedData {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

}