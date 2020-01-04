package api

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/static"
	"net/http"
)

func NewRouter(logger logrus.StdLogger, rootCA []byte) *mux.Router {
	r := mux.NewRouter()
	verifySvc := NewVerifyService(verifier.NewDefaultCfg(rootCA))
	r.HandleFunc("/verify", verifySvc.VerifyHandler).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static.Assets)))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	return r
}