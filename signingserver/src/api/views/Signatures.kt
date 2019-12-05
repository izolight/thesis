package ch.bfh.ti.hirtp1ganzg1.thesis.api.views

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.*
import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.ISignaturesHoldingService
import io.ktor.application.call
import io.ktor.http.ContentDisposition
import io.ktor.http.ContentType
import io.ktor.http.HttpHeaders
import io.ktor.http.HttpStatusCode
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.Location
import io.ktor.locations.get
import io.ktor.response.header
import io.ktor.response.respond
import io.ktor.response.respondBytes
import io.ktor.routing.Routing
import org.koin.ktor.ext.inject

@KtorExperimentalLocationsAPI
@Location(URLs.SIGNATURES)
data class SignatureRetrievalRequest(val id: String) : Validatable<SignatureRetrievalRequest> {
    override fun validate(): Validated<SignatureRetrievalRequest> =
        if (this.id.isNotEmpty()) Valid(this) else Invalid(InvalidJSONException("Empty id"))
}

@KtorExperimentalLocationsAPI
fun Routing.signature() {
    val signaturesHoldingService by inject<ISignaturesHoldingService>()

    get<SignatureRetrievalRequest> { req ->
        when (val validatedReq = req.validate()) {
            is Valid -> {
                when (val signature = signaturesHoldingService.get(req.id)) {
                    null -> call.respond(HttpStatusCode.NotFound)
                    else -> {
                        call.response.header(
                            HttpHeaders.ContentDisposition, ContentDisposition.Attachment.withParameter(
                                ContentDisposition.Parameters.FileName, "signaturefile"
                            ).toString()
                        )
                        call.respondBytes(signature, ContentType.Application.OctetStream)
                    }

                }
            }
            is Invalid -> throw validatedReq.error
        }
    }
}
