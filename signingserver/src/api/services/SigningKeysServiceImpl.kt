package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.request.url
import io.ktor.http.Url
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers
import org.bouncycastle.asn1.x500.X500Name
import org.bouncycastle.asn1.x509.*
import org.bouncycastle.cert.X509CertificateHolder
import org.bouncycastle.cert.jcajce.JcaCertStore
import org.bouncycastle.cert.jcajce.JcaX509ExtensionUtils
import org.bouncycastle.cms.CMSProcessableByteArray
import org.bouncycastle.cms.CMSSignedData
import org.bouncycastle.cms.CMSSignedDataGenerator
import org.bouncycastle.cms.jcajce.JcaSignerInfoGeneratorBuilder
import org.bouncycastle.crypto.CryptoException
import org.bouncycastle.jce.X509KeyUsage
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder
import org.bouncycastle.operator.jcajce.JcaDigestCalculatorProviderBuilder
import org.bouncycastle.pkcs.PKCS10CertificationRequest
import org.bouncycastle.pkcs.jcajce.JcaPKCS10CertificationRequestBuilder
import java.security.KeyPair
import java.security.KeyPairGenerator


class Constants {
    companion object {
        const val CRYPTO_PROVIDER = "BC"
        const val KEY_ALGORITHM = "RSA"
        const val SIGNATURE_ALGORITHM = "SHA256withRSA"
        const val RSA_KEY_BITS = 4096
    }
}

class ASN1ObjectIdentifiers {
    companion object {
        const val X509V3_CRL_DISTRIBUTION_POINT = "2.5.29.35"
    }
}

class SigningKeysServiceImpl : ISigningKeysService {
    private val keyPairGenerator = KeyPairGenerator.getInstance(Constants.KEY_ALGORITHM)
    private val keyCache = ExpireableCacheDefaultImpl<SigningKeySubjectInformation, KeyPair>()
    private val contentSignerBuilder = JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)

    init {
        keyPairGenerator.initialize(Constants.RSA_KEY_BITS)
    }

    override fun generateSigningKey(subjectInformation: SigningKeySubjectInformation): PKCS10CertificationRequest {
        val keyPair = this.keyPairGenerator.generateKeyPair()
        this.keyCache.set(subjectInformation, keyPair)

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

    override suspend fun signToPkcs7(
        subjectInformation: SigningKeySubjectInformation,
        dataToSign: ByteArray,
        signedCertificate: X509CertificateHolder
    ): CMSSignedData {
// TODO: get certificate chain, ocsp, crl
        val crl = retrieveCrl(signedCertificate)
        return CMSSignedDataGenerator().also {
            it.addSignerInfoGenerator(
                JcaSignerInfoGeneratorBuilder(
                    JcaDigestCalculatorProviderBuilder().setProvider(Constants.CRYPTO_PROVIDER).build()
                ).build(
                    JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)
                        .setProvider(Constants.CRYPTO_PROVIDER)
                        .build(this.keyCache.get(subjectInformation)!!.private),
                    signedCertificate
                )
            )
            it.addCertificates(
                JcaCertStore(listOf(signedCertificate))
            )
        }.generate(CMSProcessableByteArray(dataToSign))
    }

    @Serializable
    class CfsslCrlResponse(
        val success: Boolean,
        val result: String,
        val errors: List<CertificateAuthorityServiceImpl.ResponseMessage>,
        val messages: List<CertificateAuthorityServiceImpl.ResponseMessage>
    ) : Validatable<CfsslCrlResponse> {
        override fun validate(): Validated<CfsslCrlResponse> {
            return when {
                errors.isEmpty() and success -> Valid(this)
                else -> Invalid(InvalidDataException(errors[0].message))
            }
        }
    }

    suspend fun extractCrlUrl(signedCertificate: X509CertificateHolder): Url {
        return Url(
            withContext(Dispatchers.IO) {
                (CRLDistPoint.getInstance(
                    JcaX509ExtensionUtils.parseExtensionValue(
                        signedCertificate.getExtension(Extension.cRLDistributionPoints).extnValue.encoded
                    )
                ).distributionPoints[0].distributionPoint.name as GeneralNames).names[0].name.toString()
            }
        )
    }

    suspend fun retrieveCrl(signedCertificate: X509CertificateHolder): String {
        return when (
            val response = HttpClient {
                defaultConfig()
            }.use {
                it.get<Validatable<CfsslCrlResponse>> {
                    url(
                        extractCrlUrl(signedCertificate)
                    )
                }
            }.validate()
            ) {
            is Valid -> response.value.result
            is Invalid -> throw response.error
        }
    }
}