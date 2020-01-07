package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import Signature
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidDataException
import com.auth0.jwt.interfaces.DecodedJWT
import kotlinx.serialization.Serializable
import org.bouncycastle.cert.jcajce.JcaCertStore
import org.bouncycastle.cert.jcajce.JcaX509CertificateHolder
import org.bouncycastle.cms.CMSSignedData
import org.bouncycastle.pkcs.PKCS10CertificationRequest

@Serializable
data class SigningKeySubjectInformation(val surname: String, val givenName: String, val email: String) {
    companion object Constants {
        const val ORGANISATIONAL_UNIT = "Demo Signing Service"
        const val COUNTRY = "CH"

        fun fromIdToken(idToken: DecodedJWT): SigningKeySubjectInformation = try {
            SigningKeySubjectInformation(
                email = idToken.getClaim("email").asString()!!,
                surname = idToken.getClaim("family_name").asString()!!,
                givenName = idToken.getClaim("given_name").asString()!!
            )
        } catch (e: NullPointerException) {
            throw InvalidDataException("Required claim missing in id_token")
        }
    }

    fun toDN(): String = "CN=${surname.toUpperCase()} $givenName,OU=$ORGANISATIONAL_UNIT,DC=$COUNTRY"

}

interface ISigningKeysService {
    fun generateSigningKey(subjectInformation: SigningKeySubjectInformation): PKCS10CertificationRequest
    fun destroySigningKey(subjectInformation: SigningKeySubjectInformation)
    suspend fun signCMS(
        subjectInformation: SigningKeySubjectInformation,
        dataToSign: Signature.SignatureData,
        signedCertificate: JcaX509CertificateHolder
    ): CMSSignedData

    suspend fun fetchBundle(cert: JcaX509CertificateHolder): JcaCertStore
}


interface ICertificateAuthorityService {
    suspend fun signCSR(certificateSigningRequest: PKCS10CertificationRequest): JcaX509CertificateHolder
}


