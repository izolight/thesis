package main

import (
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/verifiy", verifier.VerifyHandler).Methods("POST")
	r.Handle("/static", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
