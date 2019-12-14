package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/config"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/static"
	"io/ioutil"
	"net/http"
)

func main() {
	// Serve static files
	cfg := config.Assets
	rootCAFile, err := cfg.Open("rootCA.pem")
	if err != nil {
		log.Fatalln(err)
	}
	rootCA, err := ioutil.ReadAll(rootCAFile)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	verifySvc := verifier.NewVerifyService(verifier.NewDefaultCfg(rootCA))
	r.HandleFunc("/verify", verifySvc.VerifyHandler).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static.Assets)))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})
	log.Fatalln(http.ListenAndServe(":8080", r))
}
