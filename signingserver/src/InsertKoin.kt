package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.*
import org.koin.dsl.module

val DIModule = module {
    single<INonceGeneratorService> { NonceGeneratorServiceDefaultImpl() }
    single<IHashesCachingService> { ExpirableHashesCachingServiceImpl() }
    single<ISignaturesHoldingService> { SignaturesHoldingServiceDefaultImpl() }
}
