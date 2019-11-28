package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

import io.ktor.client.HttpClientConfig
import io.ktor.client.features.json.JsonFeature
import io.ktor.client.features.json.serializer.KotlinxSerializer
import io.ktor.client.features.logging.LogLevel
import io.ktor.client.features.logging.Logging


fun HttpClientConfig<*>.defaultConfig() {
    install(Logging) {
        level = LogLevel.ALL
//                level = LogLevel.INFO
    }
    install(JsonFeature) {
        serializer = KotlinxSerializer()
    }


}
