package backend

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
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

var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func cookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "SESSIONID")
		session.Save(r, w)

		next.ServeHTTP(w, r)
	})
}
