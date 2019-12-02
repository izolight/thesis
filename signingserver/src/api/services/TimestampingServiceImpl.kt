package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.sha256
import io.ktor.client.HttpClient
import io.ktor.client.features.logging.LogLevel
import io.ktor.client.features.logging.Logging
import io.ktor.client.request.post
import io.ktor.client.request.url
import io.ktor.http.ContentType
import io.ktor.http.content.ByteArrayContent
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.bouncycastle.cms.CMSAlgorithm
import org.bouncycastle.tsp.TimeStampRequestGenerator
import java.io.File

class TimestampingServiceImpl : ITimestampingService {
    companion object {
        const val TSA_URL = "http://tsa.swisssign.net"
    }

    override suspend fun stamp(dataToStamp: ByteArray): ByteArray = withContext(Dispatchers.IO) {
        HttpClient {
            install(Logging) {
                level = LogLevel.HEADERS
            }
        }.use {
            it.post<ByteArray> {
                url(TSA_URL)
                body = ByteArrayContent(
                    TimeStampRequestGenerator().also { gen ->
                        gen.setCertReq(true)
                    }.generate(
                        CMSAlgorithm.SHA256,
                        sha256(dataToStamp)
                    ).encoded.also { tsq ->
                        File("/tmp/tsa_req").writeBytes(tsq)
                    },
                    contentType = ContentType("application", "timestamp-query")
                )
            }
        }.also {
            File("/tmp/tsa_resp").writeBytes(it)
        }
    }
}