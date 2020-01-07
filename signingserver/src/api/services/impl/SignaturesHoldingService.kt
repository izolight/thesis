package ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def.ISignaturesHoldingService
import java.security.SecureRandom
import kotlin.math.absoluteValue


class SignaturesHoldingServiceDefaultImpl : ICacheDefaultImpl<String, ByteArray>(), ISignaturesHoldingService {
    private val random = SecureRandom()

    override fun generateId(): String = random.nextLong().absoluteValue.toString(16)
}
