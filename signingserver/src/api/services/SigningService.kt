package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import org.bouncycastle.pkcs.PKCS10CertificationRequest
import javax.security.cert.X509Certificate

data class SigningKeySubjectInformation(val surname: String, val givenName: String, val email: String) {
    companion object Constants {
        val ORGANISATIONAL_UNIT = "Demo Signing Service"
        val COUNTRY = "CH"
    }

    fun toDN(): String {
        return "CN=${surname.toUpperCase()} $givenName,OU=$ORGANISATIONAL_UNIT,DC=$COUNTRY"
    }
}

interface SigningKeysService {
    fun generateSigningKey(subjectInformation: SigningKeySubjectInformation): PKCS10CertificationRequest
}


interface CertificateAuthorityService {
    fun signCSR(certificateSigningRequest: PKCS10CertificationRequest): X509Certificate
}


