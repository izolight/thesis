package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.*
import io.ktor.util.KtorExperimentalAPI
import kotlinx.coroutines.runBlocking
import org.koin.dsl.module

@KtorExperimentalAPI
val DIModule = module {
    single<INonceGeneratorService> { NonceGeneratorServiceDefaultImpl() }
    single<ISecretService> { SecretServiceDefaultImpl() }
    single<ISignaturesHoldingService> { SignaturesHoldingServiceDefaultImpl() }
    single<IOIDCService> { runBlocking { OurDemoOIDCService() } }
    single<ISigningKeysService> { SigningKeysServiceImpl() }
    single<ICertificateAuthorityService> { CertificateAuthorityServiceImpl() }
}
