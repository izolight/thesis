package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import java.security.SecureRandom

interface ISecretService {
    fun getSecret(): Long
    fun getHmacKey(nonce: Long): ByteArray
}


class SecretServiceDefaultImpl : ISecretService {
    private val secret = SecureRandom().nextLong()
    override fun getSecret(): Long {
        return this.secret
    }

    override fun getHmacKey(nonce: Long): ByteArray {
        return (
                this.getSecret().toString() + nonce.toString().toByteArray()
                ).toByteArray()
    }
}