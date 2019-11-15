package ch.bfh.ti.hirtp1ganzg1.thesis

import io.ktor.http.Url

class Constants {
    companion object {
        val OIDC_CLIENT_ID = "493445436490-6prvgh3d4ubac679519mhg5rlhokqni2.apps.googleusercontent.com"
        val OIDC_CLIENT_SECRET = "LjFQb7iPQzfUFmna099gr0vm"
        val OIDC_REDIRECT_URI = Url("https://thesis.izolight.xyz/oidc-redirect")
        val OIDC_SCOPES = listOf("openid", "email")
        val OIDC_RESPONSE_TYPE = "id_token"
    }
}