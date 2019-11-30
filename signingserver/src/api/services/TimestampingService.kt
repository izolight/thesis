package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

interface ITimestampingService {
    suspend fun stamp(dataToStamp: ByteArray): ByteArray
}