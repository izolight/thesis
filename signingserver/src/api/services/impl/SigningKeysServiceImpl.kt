package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl

import Signature
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def.ISigningKeysService
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import io.ktor.client.HttpClient
import io.ktor.client.engine.apache.Apache
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.request.url
import io.ktor.http.ContentType
import io.ktor.http.Url
import io.ktor.http.content.ByteArrayContent
import io.ktor.http.contentType
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import org.bouncycastle.asn1.ASN1Encoding
import org.bouncycastle.asn1.ASN1OutputStream
import org.bouncycastle.asn1.DEROctetString
import org.bouncycastle.asn1.ocsp.OCSPObjectIdentifiers
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers
import org.bouncycastle.asn1.x500.X500Name
import org.bouncycastle.asn1.x509.*
import org.bouncycastle.cert.X509CertificateHolder
import org.bouncycastle.cert.jcajce.JcaCertStore
import org.bouncycastle.cert.jcajce.JcaX509CRLHolder
import org.bouncycastle.cert.jcajce.JcaX509CertificateHolder
import org.bouncycastle.cert.jcajce.JcaX509ExtensionUtils
import org.bouncycastle.cert.ocsp.CertificateID
import org.bouncycastle.cert.ocsp.OCSPReqBuilder
import org.bouncycastle.cert.ocsp.OCSPResp
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
import org.slf4j.LoggerFactory
import java.io.ByteArrayInputStream
import java.io.ByteArrayOutputStream
import java.security.KeyPair
import java.security.KeyPairGenerator
import java.security.SecureRandom
import java.security.cert.CertificateFactory
import java.security.cert.X509CRL
import java.util.*


class Constants {
    companion object {
        const val KEY_ALGORITHM = "RSA"
        const val SIGNATURE_ALGORITHM = "SHA256withRSA"
        const val RSA_KEY_BITS = 4096
    }
}

class SigningKeysServiceImpl : ISigningKeysService {
    private val keyPairGenerator = KeyPairGenerator.getInstance(Constants.KEY_ALGORITHM)
    private val keyCache = HashMap<ISigningKeysService.SigningKeySubjectInformation, KeyPair>()
    private val contentSignerBuilder = JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)
    private val secureRandom = SecureRandom()
    private val logger = LoggerFactory.getLogger(SigningKeysServiceImpl::class.java)

    init {
        keyPairGenerator.initialize(Constants.RSA_KEY_BITS)
    }

    override fun generateSigningKey(subjectInformation: ISigningKeysService.SigningKeySubjectInformation): PKCS10CertificationRequest {
        val keyPair = this.keyPairGenerator.generateKeyPair()
        this.keyCache[subjectInformation] = keyPair

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

    override fun destroySigningKey(subjectInformation: ISigningKeysService.SigningKeySubjectInformation) {
        this.keyCache.remove(subjectInformation)
    }

    private fun extractIssuerCertificate(
        cert: JcaX509CertificateHolder,
        bundle: JcaCertStore
    ): JcaX509CertificateHolder = bundle.getMatches(null)
        .filterIsInstance<JcaX509CertificateHolder>()
        .filter {
            it.subject == cert.issuer
        }[0]

    override suspend fun signCMS(
        subjectInformation: ISigningKeysService.SigningKeySubjectInformation,
        dataToSign: Signature.SignatureData,
        signedCertificate: JcaX509CertificateHolder
    ): CMSSignedData = CMSSignedDataGenerator().also {
        withContext(Dispatchers.IO) {
            val futureBundle = async { fetchBundle(signedCertificate) }
            val crlCert = async { retrieveCrl(signedCertificate) }
            val bundle = futureBundle.await()
            val issuerCert = extractIssuerCertificate(signedCertificate, bundle)
            val rootCert = extractIssuerCertificate(issuerCert, bundle)
            val ocspCert = async { retrieveOcsp(signedCertificate, issuerCert) }
            val ocspIssuer = async { retrieveOcsp(issuerCert, rootCert) }
            it.addSignerInfoGenerator(
                JcaSignerInfoGeneratorBuilder(
                    JcaDigestCalculatorProviderBuilder().build()
                ).build(
                    JcaContentSignerBuilder(Constants.SIGNATURE_ALGORITHM)
                        .build(keyCache[subjectInformation]!!.private),
                    signedCertificate
                )
            )
            crlCert.await().ifPresent { crl ->
                it.addCRL(crl)
            }
            it.addOtherRevocationInfo(
                OCSPObjectIdentifiers.id_pkix_ocsp_basic,
                ocspCert.await().toASN1Structure()
            )
            it.addOtherRevocationInfo(
                OCSPObjectIdentifiers.id_pkix_ocsp_basic,
                ocspIssuer.await().toASN1Structure()
            )
            it.addCertificates(bundle)
        }
    }.generate(CMSProcessableByteArray(dataToSign.toByteArray()), true)


    @Suppress("MemberVisibilityCanBePrivate", "unused")
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

    private suspend fun extractCrlUrl(signedCertificate: X509CertificateHolder): List<Url> =
        withContext(Dispatchers.IO) {
            CRLDistPoint.getInstance(
                JcaX509ExtensionUtils.parseExtensionValue(
                    signedCertificate.getExtension(Extension.cRLDistributionPoints).extnValue.encoded
                )
            ).distributionPoints.mapNotNull {
                try {
                    Url((it.distributionPoint.name as GeneralNames).names[0].name.toString())
                } catch (e: Exception) {
                    logger.error("Unable to parse CRL distribution point: %s", e)
                    null
                }
            }
        }

    private suspend fun retrieveCrl(signedCertificate: X509CertificateHolder): Optional<JcaX509CRLHolder> {
        HttpClient(Apache) { defaultConfig() }.use { client ->
            extractCrlUrl(signedCertificate).forEach { url ->
                try {
                    when (val response = client.get<CfsslCrlResponse> { url(url) }.validate()) {
                        is Valid -> return Optional.of(
                            JcaX509CRLHolder(
                                CertificateFactory.getInstance("X.509").generateCRL(
                                    ByteArrayInputStream(
                                        Base64.getDecoder().decode(response.value.result)
                                    )
                                ) as X509CRL
                            )
                        )
                        is Invalid -> throw response.error

                    }
                } catch (e: Exception) {
                    logger.info("Failed to download CRL from %s, trying next", url)
                }
            }
        }
        logger.warn("Unable to reach any CRL distribution points")
        return Optional.empty()
    }

    private fun constructOcspRequest(
        signedCertificate: X509CertificateHolder,
        issuer: X509CertificateHolder
    ) = OCSPReqBuilder()
        .addRequest(
            CertificateID(
                JcaDigestCalculatorProviderBuilder()
                    .build()
                    .get(CertificateID.HASH_SHA1),
                issuer,
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

    private suspend fun retrieveOcsp(
        signedCertificate: X509CertificateHolder,
        issuer: X509CertificateHolder
    ) = withContext(Dispatchers.IO) {
        OCSPResp(
            HttpClient(Apache) {
                defaultConfig()
            }.use {
                it.post<ByteArray> {
                    url(extractOcspUrl(signedCertificate))
                    body = ByteArrayContent(
                        bytes = constructOcspRequest(
                            signedCertificate,
                            issuer
                        ).encoded,
                        contentType = ContentType(
                            "application", "ocsp-request"
                        )
                    )
                }
            }
        )
    }

    override suspend fun fetchBundle(cert: JcaX509CertificateHolder) =
        JcaCertStore(
            HttpClient(Apache) {
                defaultConfig()
            }.use {
                it.post<CertificateAuthorityServiceImpl.CfsslBundleResponse> {
                    url(CertificateAuthorityServiceImpl.CA_BUNDLE_URL)
                    contentType(ContentType.Application.Json)
                    body = CertificateAuthorityServiceImpl.CfsslBundleRequest(certificate = certificateToPem(
                        cert
                    )
                    )
                }.result.allCerts()
            })
}

fun OCSPResp.toDER(): ByteArray = ByteArrayOutputStream().also {
    ASN1OutputStream.create(it, ASN1Encoding.DER).also { asn1outputStream ->
        asn1outputStream.writeObject(this.toASN1Structure())
        asn1outputStream.close()
    }
}.toByteArray()
