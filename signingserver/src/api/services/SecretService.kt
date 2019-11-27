package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

interface ISecretService {
    fun getSecret(): ByteArray
}


class SecretServiceDefaultImpl : ISecretService {
    private val secret = ByteArray(32).also { SecureRandom().nextBytes(it) }
    override fun getSecret(): ByteArray {
        return this.secret.copyOf()
    }

}