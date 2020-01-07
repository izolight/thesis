package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def


interface ISecretService {
    fun hkdf(salt: ByteArray, length: Int = 64): ByteArray
}
