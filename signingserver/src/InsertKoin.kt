package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.*
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.client.features.json.JsonFeature
import io.ktor.client.features.json.serializer.KotlinxSerializer
import io.ktor.client.features.logging.LogLevel
import io.ktor.client.features.logging.Logging
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
    single {
        HttpClient(CIO) {
            install(Logging) {
                level = LogLevel.ALL
//                level = LogLevel.INFO
            }
            install(JsonFeature) {
                serializer = KotlinxSerializer()
            }
        }
    }
}
