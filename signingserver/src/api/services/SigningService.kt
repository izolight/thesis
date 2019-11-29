package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import Signature
import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.Either
import com.auth0.jwt.interfaces.DecodedJWT
import org.bouncycastle.cms.CMSSignedData
import org.bouncycastle.pkcs.PKCS10CertificationRequest
import org.slf4j.LoggerFactory

data class SigningKeySubjectInformation(val surname: String, val givenName: String, val email: String) {
    companion object Constants {
        private val logger = LoggerFactory.getLogger(this.javaClass)
        const val ORGANISATIONAL_UNIT = "Demo Signing Service"
        const val COUNTRY = "CH"

        fun fromIdToken(idToken: DecodedJWT): Either<SigningKeySubjectInformation> {
            return try {
                Either.Success(
                    SigningKeySubjectInformation(
                        email = idToken.getClaim("email").asString()!!,
                        surname = idToken.getClaim("family_name").asString()!!,
                        givenName = idToken.getClaim("given_name").asString()!!
                    )
                )
            } catch (e: NullPointerException) {
                logger.error("Missing required claim", e)
                Either.Error("Missing required claim", e)
            }
        }
    }

    fun toDN(): String {
        return "CN=${surname.toUpperCase()} $givenName,OU=$ORGANISATIONAL_UNIT,DC=$COUNTRY"
    }

}

interface ISigningKeysService {
    fun generateSigningKey(subjectInformation: SigningKeySubjectInformation): PKCS10CertificationRequest
    fun signToPkcs7(signatureData: Signature.SignatureData): CMSSignedData
}


interface ICertificateAuthorityService {
    fun signCSR(certificateSigningRequest: PKCS10CertificationRequest): java.security.cert.X509Certificate
}


