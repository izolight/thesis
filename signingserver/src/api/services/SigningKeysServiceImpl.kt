package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.request.url
import io.ktor.http.ContentType
import io.ktor.http.Url
import io.ktor.http.content.ByteArrayContent
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import org.bouncycastle.asn1.DEROctetString
import org.bouncycastle.asn1.ocsp.OCSPObjectIdentifiers
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers
import org.bouncycastle.asn1.x500.X500Name
import org.bouncycastle.asn1.x509.*
import org.bouncycastle.cert.X509CertificateHolder
import org.bouncycastle.cert.jcajce.JcaCertStore
import org.bouncycastle.cert.jcajce.JcaX509CRLHolder
import org.bouncycastle.cert.jcajce.JcaX509ExtensionUtils
import org.bouncycastle.cert.ocsp.CertificateID
import org.bouncycastle.cert.ocsp.OCSPReqBuilder
import org.bouncycastle.cert.ocsp.OCSPResp
import org.bouncycastle.cms.CMSProcessableByteArray
import org.bouncycastle.cms.CMSSignedDataGenerator
import org.bouncycastle.cms.jcajce.JcaSignerInfoGeneratorBuilder
import org.bouncycastle.crypto.CryptoException
import org.bouncycastle.jce.X509KeyUsage
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder
import org.bouncycastle.operator.jcajce.JcaDigestCalculatorProviderBuilder
import org.bouncycastle.pkcs.PKCS10CertificationRequest
import org.bouncycastle.pkcs.jcajce.JcaPKCS10CertificationRequestBuilder
import java.io.ByteArrayInputStream
import java.io.File
import java.security.KeyPair
import java.security.KeyPairGenerator
import java.security.SecureRandom
import java.security.cert.CertificateFactory
import java.security.cert.X509CRL
import java.util.*


class Constants {
    companion object {
        const val CRYPTO_PROVIDER = "BC"
        const val KEY_ALGORITHM = "RSA"
        const val SIGNATURE_ALGORITHM = "SHA256withRSA"
        const val RSA_KEY_BITS = 4096
    }
}

class SigningKeysServiceImpl : ISigningKeysService {
    private val keyPairGenerator = KeyPairGenerator.getInstance(Constants.KEY_ALGORITHM)
    private val keyCache = ExpireableCacheDefaultImpl<SigningKeySubjectInformation, KeyPair>()
    private val contentSignerBuilder = JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)
    private val secureRandom = SecureRandom()

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

    // TODO: get and insert certificate chain
    override suspend fun signToPkcs7(
        subjectInformation: SigningKeySubjectInformation,
        dataToSign: ByteArray,
        signedCertificate: X509CertificateHolder,
        bundle: JcaCertStore
    ) = CMSSignedDataGenerator().also {
        it.addSignerInfoGenerator(
            JcaSignerInfoGeneratorBuilder(
                JcaDigestCalculatorProviderBuilder()
                    .build()
            ).build(
                JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)
                    .build(this.keyCache.get(subjectInformation)!!.private),
                signedCertificate
            )
        )
        it.addCertificates(
            bundle
        )
        it.addCRL(
            retrieveCrl(signedCertificate)
        )
        it.addOtherRevocationInfo(
            OCSPObjectIdentifiers.id_pkix_ocsp,
            retrieveOcsp(signedCertificate).toASN1Structure()
        )
    }.generate(CMSProcessableByteArray(dataToSign))

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

    private fun extractOcspUrl(signedCertificate: X509CertificateHolder) = Url(
        AuthorityInformationAccess.fromExtensions(
            signedCertificate.extensions
        ).accessDescriptions[0].accessLocation.name.toString()
    )

    // TODO this is not optimal: only one crl url is extracted, and there is no error handling
    private suspend fun extractCrlUrl(signedCertificate: X509CertificateHolder) = Url(
        withContext(Dispatchers.IO) {
            (CRLDistPoint.getInstance(
                JcaX509ExtensionUtils.parseExtensionValue(
                    signedCertificate.getExtension(Extension.cRLDistributionPoints).extnValue.encoded
                )
            ).distributionPoints[0].distributionPoint.name as GeneralNames).names[0].name.toString()
        }
    )

    suspend fun retrieveCrl(signedCertificate: X509CertificateHolder) = when (
        val response = HttpClient {
            defaultConfig()
        }.use {
            it.get<CfsslCrlResponse> {
                url(extractCrlUrl(signedCertificate))
            }
        }.validate()) {
        is Valid -> JcaX509CRLHolder(
            CertificateFactory.getInstance("X.509").generateCRL(
                ByteArrayInputStream(
                    Base64.getDecoder().decode(response.value.result)
                )
            ) as X509CRL
        )
        is Invalid -> throw response.error
    }

    private fun constructOcspRequest(signedCertificate: X509CertificateHolder) = OCSPReqBuilder()
        .addRequest(
            CertificateID(
                JcaDigestCalculatorProviderBuilder()
                    .build()
                    .get(CertificateID.HASH_SHA1),
                signedCertificate,
                signedCertificate.serialNumber
            )
        )
        .setRequestExtensions(
            Extensions(
                Extension(
                    OCSPObjectIdentifiers.id_pkix_ocsp_nonce,
                    false,
                    DEROctetString(
                        ByteArray(32).also {
                            secureRandom.nextBytes(it)
                        }
                    )
                )
            )
        ).build()

    suspend fun retrieveOcsp(signedCertificate: X509CertificateHolder) = withContext(Dispatchers.IO) {
        OCSPResp(
            HttpClient {
                defaultConfig()
            }.use {
                it.post<ByteArray> {
                    url(extractOcspUrl(signedCertificate))
                    body = ByteArrayContent(
                        bytes = constructOcspRequest(signedCertificate).encoded.also { bytes ->
                            File("/tmp/ocsp_req").writeBytes(bytes)
                        },
                        contentType = ContentType(
                            "application", "ocsp-request"
                        )
                    )
                }
            }
        )
    }
}
