package main

import (
	"gitlab.ti.bfh.ch/hirtp1/thesis/backend"
	"log"
	"net/http"
	"time"
)

func main() {
	c, err := backend.NewOIDCClient(
		"https://accounts.google.com",
		"ID",
		"Secret",
		"https://localhost:8080",
	)
	if err != nil {
		log.Fatal(err)
	}

	r := backend.NewRouter(c)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
