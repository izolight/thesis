package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hexStringToByteArray
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.hmacSha256
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.request.url
import io.ktor.http.ContentType
import io.ktor.http.contentType
import kotlinx.io.StringWriter
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration
import org.bouncycastle.cert.jcajce.JcaCertStore
import org.bouncycastle.cert.jcajce.JcaX509CertificateHolder
import org.bouncycastle.openssl.MiscPEMGenerator
import org.bouncycastle.openssl.jcajce.JcaPEMWriter
import org.bouncycastle.pkcs.PKCS10CertificationRequest
import org.bouncycastle.util.io.pem.PemWriter
import java.io.ByteArrayInputStream
import java.security.cert.CertificateFactory
import java.security.cert.X509Certificate
import java.util.*

class CertificateAuthorityServiceImpl : ICertificateAuthorityService {
    companion object {
        private const val HMAC_KEY = "7BAFD191E2631D4505F612C7D6B2010A"
        const val CA_URL = "https://intermediate-ca.thesis.izolight.xyz"
        const val CA_SIGN_URL = "${CA_URL}/api/v1/cfssl/authsign"
        const val CA_BUNDLE_URL = "${CA_URL}/api/v1/cfssl/bundle"
        private val json = Json(JsonConfiguration.Stable)
    }

    @Serializable
    data class Request(val token: String, val request: String)

    @Serializable
    data class CertificateRequest(val certificate_request: String, val profile: String = "signingService") :
        Validatable<CertificateRequest> {
        override fun validate(): Validated<CertificateRequest> {
            if (this.certificate_request.isNotEmpty()) {
                return Valid(this)
            }
            return Invalid(IllegalArgumentException("Certificate is empty"))
        }
    }

    @Serializable
    data class ResponseMessage(val code: Int, val message: String)

    @Serializable
    data class CfsslSignResult(val certificate: String)

    @Serializable
    data class CfsslResponse(
        val success: Boolean,
        val result: CfsslSignResult,
        val errors: List<ResponseMessage>,
        val messages: List<ResponseMessage>
    ) : Validatable<CfsslResponse> {
        override fun validate(): Validated<CfsslResponse> {
            return when {
                errors.isEmpty() and success -> Valid(this)
                else -> Invalid(InvalidDataException(errors[0].message))
            }
        }
    }

    @Serializable
    data class CfsslBundleStatus(
        val code: Int,
        val expiring_SKIs: List<String>,
        val messages: List<String>,
        val rebundled: Boolean,
        val untrusted_root_stores: List<String>
    )

    @Serializable
    data class CfsslBundle(
        val bundle: String,
        val crl_support: Boolean,
        val crt: String,
        val expires: String,
        val hostnames: List<String>,
        val issuer: String,
        val key: String,
        val key_size: Int,
        val key_type: String,
        val leaf_expires: String,
        val ocsp: List<String>,
        val ocsp_support: Boolean,
        val root: String,
        val signature: String,
        val status: CfsslBundleStatus,
        val subject: String
    ) {
        fun splitBundleIntoPems() =
            this.bundle.splitWithDelimiters("-----END CERTIFICATE-----").map {
                it.removePrefix("\n")
            }

        fun allPems() = splitBundleIntoPems().toMutableList().also { it.add(this.root) }.toList()

        fun allCerts() = allPems().map { s -> pemToCertificate(s) }
    }

    @Serializable
    data class CfsslBundleResponse(
        val success: Boolean,
        val result: CfsslBundle,
        val errors: List<ResponseMessage>,
        val messages: List<ResponseMessage>
    )

    @Serializable
    data class CfsslBundleRequest(
        val certificate: String
    )

    private fun authenticateCertificateRequest(request: Valid<CertificateRequest>) = with(
        json.toJson(
            CertificateRequest.serializer(),
            request.value
        ).toString().toByteArray(Charsets.UTF_8)
    ) {
        Request(
            request = Base64.getEncoder().encodeToString(this),
            token = Base64.getEncoder().encodeToString(
                hmacSha256(
                    hexStringToByteArray(HMAC_KEY),
                    this
                )
            )
        )
    }

    override suspend fun signCSR(certificateSigningRequest: PKCS10CertificationRequest) = when (
        val validatedResponse = HttpClient { defaultConfig() }.use {
            it.post<CfsslResponse> {
                url(CA_SIGN_URL)
                contentType(ContentType.Application.Json)
                body = authenticateCertificateRequest(
                    when (val req = CertificateRequest(
                        certificate_request = StringWriter().also { writer ->
                            PemWriter(writer).also { p ->
                                MiscPEMGenerator(certificateSigningRequest).also { generator ->
                                    p.writeObject(generator)
                                }
                                p.close()
                            }
                        }.toString()
                    ).validate()) {
                        is Valid -> req
                        is Invalid -> throw req.error
                    }
                )
            }
        }.validate()
        ) {
        is Valid -> pemToCertificate(validatedResponse.value.result.certificate)
        is Invalid -> throw validatedResponse.error
    }


    override suspend fun fetchBundle(cert: JcaX509CertificateHolder) = JcaCertStore(
        HttpClient {
            defaultConfig()
        }.use {
            it.post<CfsslBundleResponse> {
                url(CA_BUNDLE_URL)
                contentType(ContentType.Application.Json)
                body = CfsslBundleRequest(certificate = certificateToPem(cert))
            }.result.allCerts()
        })
}

fun String.splitWithDelimiters(delimiter: String): List<String> =
    this.split(delimiter).filter { s -> s.isNotEmpty() }.map { s -> s.plus(delimiter) }


fun pemToCertificate(pem: String) = JcaX509CertificateHolder(
    CertificateFactory.getInstance("X.509")
        .generateCertificate(
            ByteArrayInputStream(
                pem.toByteArray(Charsets.UTF_8)
            )
        ) as X509Certificate
)

fun certificateToPem(cert: JcaX509CertificateHolder) = StringWriter().also {
    JcaPEMWriter(it).also { pemWriter ->
        pemWriter.writeObject(cert)
        pemWriter.close()
    }
}.toString()
