package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def

import org.bouncycastle.cert.jcajce.JcaX509CertificateHolder
import org.bouncycastle.pkcs.PKCS10CertificationRequest

interface ICertificateAuthorityService {
    suspend fun signCSR(certificateSigningRequest: PKCS10CertificationRequest): JcaX509CertificateHolder
}
