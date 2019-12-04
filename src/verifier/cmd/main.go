package main

import (
	"github.com/gorilla/mux"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/verify", verifier.VerifyHandler).Methods("POST")
	r.Handle("/static", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
