package backend

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(client *OIDCClient) *mux.Router {
	dir := "./static"
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/api/hashes", uploadHashHandler).Methods("POST")
	r.HandleFunc("/oauth-redirect", client.oidcRedirectHandler)
	r.Use(cookieMiddleware)

	return r
}
