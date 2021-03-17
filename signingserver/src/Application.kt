package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.ApiError
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidRequestException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.postHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.sign
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.signature
import io.ktor.application.*
import io.ktor.features.*
import io.ktor.http.*
import io.ktor.http.content.*
import io.ktor.locations.*
import io.ktor.request.*
import io.ktor.response.*
import io.ktor.routing.*
import io.ktor.serialization.*
import io.ktor.util.*
import kotlinx.serialization.json.Json
import org.koin.ktor.ext.Koin
import org.slf4j.event.Level

fun main(args: Array<String>) = io.ktor.server.jetty.EngineMain.main(args)

//fun main(args: Array<String>): Unit {
//    Security.addProvider(BouncyCastleProvider())
//    io.ktor.server.netty.EngineMain.main(args)
//}

@KtorExperimentalLocationsAPI
@KtorExperimentalAPI
fun Application.module() {

    install(Compression) {
        gzip {
            priority = 1.0
        }
        deflate {
            priority = 10.0
            minimumSize(1024) // condition
        }
    }

    install(CallLogging) {
        level = Level.INFO
        filter { call -> call.request.path().startsWith("/") }
    }

    install(Koin) {
        modules(DIModule)
    }

    install(DefaultHeaders) {
        header("X-Engine", "Ktor") // will send this header with each response
    }

    install(ContentNegotiation) {
        json(
            Json {
                prettyPrint = true
            },
            contentType = ContentType.Application.Json
        )
    }

    install(StatusPages) {
        exception<InvalidRequestException> { exception ->
            call.respond(
                HttpStatusCode.BadRequest,
                ApiError(message = exception.message ?: "Unspecified error")
            )
        }

        exception<Throwable> { exception ->
            call.application.environment.log.error(
                "Unhandled exception",
                exception
            )
            call.respond(
                HttpStatusCode.InternalServerError,
                ApiError(message = exception.message ?: "Unspecified error")
            )
        }
    }

    install(Locations)

    routing {
        //        trace { application.log.trace(it.buildText()) }

        root()
        static()
        callback()
        postHashes()
        sign()
        signature()

    }
}

fun Routing.root() {
    get("/") {
        call.respondRedirect("/static/index.html")
    }
}

fun Routing.static() {
    static("/static") {
        //        files("resources/static")
        resources("static")
    }

}

fun Routing.callback() {
    get("/callback") {
        call.respondRedirect("/static/callback.html")
    }
}
