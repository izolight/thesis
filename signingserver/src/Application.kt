package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.ApiError
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidRequestException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.postHashes
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.sign
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.signature
import io.ktor.application.Application
import io.ktor.application.call
import io.ktor.application.install
import io.ktor.features.*
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.content.files
import io.ktor.http.content.static
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.Locations
import io.ktor.request.path
import io.ktor.response.respond
import io.ktor.response.respondRedirect
import io.ktor.routing.Routing
import io.ktor.routing.get
import io.ktor.routing.routing
import io.ktor.serialization.DefaultJsonConfiguration
import io.ktor.serialization.serialization
import io.ktor.util.KtorExperimentalAPI
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
@Suppress("unused") // Referenced in application.conf
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
        //        jackson {
//            configure(SerializationFeature.INDENT_OUTPUT, true)
//            setDefaultPrettyPrinter(DefaultPrettyPrinter().apply {
//                indentArraysWith(DefaultPrettyPrinter.FixedSpaceIndenter.instance)
//                indentObjectsWith(DefaultIndenter("  ", "\n"))
//            })
//            registerModule(KotlinModule())
//        }
        serialization(
            contentType = ContentType.Application.Json,
            json = Json(
                DefaultJsonConfiguration.copy(
                    prettyPrint = true
                )
            )

        )
    }

    install(StatusPages) {
        exception<InvalidRequestException> { exception ->
            call.respond(
                HttpStatusCode.BadRequest,
                ApiError(message = "Invalid request: ${exception.message}")
            )
        }

        exception<Throwable> { exception ->
            call.application.environment.log.error(
                "Unhandled exception",
                exception
            )
            call.respond(
                HttpStatusCode.InternalServerError,
                ApiError(message = "Unexpected error: ${exception.message ?: "Unknown"}")
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
        files("resources/static")
    }

}

fun Routing.callback() {
    get("/callback") {
        call.respondRedirect("/static/callback.html")
    }
}
