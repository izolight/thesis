package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import org.bouncycastle.pkcs.PKCS10CertificationRequest
import javax.security.cert.X509Certificate

class CertificateAuthorityServiceImpl : CertificateAuthorityService {
    override fun signCSR(certificateSigningRequest: PKCS10CertificationRequest): X509Certificate {
        TODO()
    }
}