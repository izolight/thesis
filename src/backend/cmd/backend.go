package main

import (
	backend "gitlab.ti.bfh.ch/hirtp1/thesis"
	"log"
	"net/http"
	"time"
)

func main() {
	r := backend.NewRouter()

	srv := &http.Server{
		Handler: r,
		Addr: "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
