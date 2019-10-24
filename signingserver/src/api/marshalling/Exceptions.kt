package ch.bfh.ti.hirtp1ganzg1.thesis.api.marshalling

//import com.fasterxml.jackson.core.JsonProcessingException

open class InvalidRequestException(message: String) : RuntimeException(message)


class InvalidJSONException(message: String) : InvalidRequestException(message)

