package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def

interface ISignaturesHoldingService :
    ICache<String, ByteArray> {
    fun generateId(): String
}
