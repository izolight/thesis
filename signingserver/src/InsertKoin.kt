package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.def.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl.*
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.runBlocking
import org.koin.dsl.module
import org.slf4j.Logger
import org.slf4j.LoggerFactory

@KtorExperimentalAPI
val DIModule = module {
    single<INonceGeneratorService> { NonceGeneratorServiceDefaultImpl() }
    single<ISecretService> { SecretServiceDefaultImpl() }
    single<ISignaturesHoldingService> { SignaturesHoldingServiceDefaultImpl() }
    single<IOIDCService> { runBlocking { OurDemoOIDCService() } }
    single<ISigningKeysService> { SigningKeysServiceImpl() }
    single<ICertificateAuthorityService> { CertificateAuthorityServiceImpl() }
    single<ITimestampingService> { TimestampingServiceImpl() }
    single<Logger> { LoggerFactory.getLogger("DEMO") }
}
