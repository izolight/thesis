package backend

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NewRouter(client *OIDCClient, logger *logrus.Logger) *mux.Router {
	logger.SetLevel(logrus.InfoLevel)
	dir := "./static"
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(dir))))
	r.Handle("/api/hashes", handlers.CombinedLoggingHandler(logger.Writer(), http.HandlerFunc(uploadHashHandler))).Methods("POST")
	r.HandleFunc("/oauth-redirect", client.oidcRedirectHandler)
	r.Use(cookieMiddleware)

	return r
}
