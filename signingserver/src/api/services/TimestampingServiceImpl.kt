package ch.bfh.ti.hirtp1ganzg1.thesis.api.services

import ch.bfh.ti.hirtp1ganzg1.thesis.api.utils.defaultConfig
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.request.url
import io.ktor.client.response.HttpResponse
import io.ktor.http.ContentType
import io.ktor.http.contentType
import org.bouncycastle.asn1.ASN1ObjectIdentifier
import org.bouncycastle.tsp.TimeStampRequestGenerator

class TimestampingServiceImpl : ITimestampingService {
    companion object {
        const val TSA_URL = "http://tsa.swisssign.net"
    }
    override suspend fun stamp(data: ByteArray): Any {
        val encodedRequest = TimeStampRequestGenerator().generate(
            ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.1"),
            data
        ).encoded
        val response = HttpClient { defaultConfig() }.use {
            it.post<HttpResponse> {
                url(TSA_URL)
                contentType(ContentType("application", "timestamp-query"))
                body = encodedRequest
            }
        }
        // TODO("run this and see what it returns")
        println()
    }
}