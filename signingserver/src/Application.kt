package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.ApiError
import ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling.InvalidJSONException
import ch.bfh.ti.hirtp1ganzg1.thesis.api.views.postHashes
import io.ktor.application.Application
import io.ktor.application.call
import io.ktor.application.install
import io.ktor.features.*
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.content.resources
import io.ktor.http.content.static
import io.ktor.locations.KtorExperimentalLocationsAPI
import io.ktor.locations.Location
import io.ktor.locations.Locations
import io.ktor.locations.get
import io.ktor.request.path
import io.ktor.response.respond
import io.ktor.response.respondText
import io.ktor.routing.Routing
import io.ktor.routing.get
import io.ktor.routing.routing
import io.ktor.serialization.DefaultJsonConfiguration
import io.ktor.serialization.serialization
import kotlinx.serialization.json.Json
import org.koin.ktor.ext.Koin
import org.slf4j.event.Level

fun main(args: Array<String>): Unit = io.ktor.server.netty.EngineMain.main(args)

@KtorExperimentalLocationsAPI
@Suppress("unused") // Referenced in application.conf
fun Application.module() {
    install(Locations) {
    }

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
        exception<InvalidJSONException> { exception ->
            call.respond(
                HttpStatusCode.BadRequest,
                ApiError("Invalid JSON: ${exception.message}")
            )
        }

        exception<Throwable> { exception ->
            call.respond(
                HttpStatusCode.InternalServerError,
                ApiError("Unexpected error: ${exception.message ?: "Unknown"}")
            )
        }
    }

//    val client = HttpClient(CIO) {
//        install(JsonFeature) {
//            serializer = GsonSerializer()
//        }
//        install(Logging) {
//            level = LogLevel.HEADERS
//        }
//    }
//    runBlocking {
    // Sample for making a HTTP Client request
    /*
    val message = client.post<JsonSampleClass> {
        url("http://127.0.0.1:8080/path/to/endpoint")
        contentType(ContentType.Application.Json)
        body = JsonSampleClass(hello = "world")
    }
    */
//    }

    routing {
        // Static feature. Try to access `/static/ktor_logo.svg`
        static("/static") {
            resources("static")
        }

        root()
        postHashes()

        get<MyLocation> {
            call.respondText("Location: name=${it.name}, arg1=${it.arg1}, arg2=${it.arg2}")
        }
    }
}

fun Routing.root() {
    get("/") {
        call.respondText("lol generics", contentType = ContentType.Text.Plain)
    }
}


@KtorExperimentalLocationsAPI
@Location("/location/{name}")
class MyLocation(val name: String, val arg1: Int = 42, val arg2: String = "default")

