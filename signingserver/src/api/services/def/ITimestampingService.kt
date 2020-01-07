package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def

interface ITimestampingService {
    suspend fun stamp(dataToStamp: ByteArray): ByteArray
}