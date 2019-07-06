package main

import (
	"gitlab.ti.bfh.ch/hirtp1/thesis/pkg"
	"log"
	"net/http"
	"time"
)

func main() {
	c, err := pkg.NewOIDCClient(
		"https://accounts.google.com",
		"ID",
		"Secret",
		"https://localhost:8080",
	)
	if err != nil {
		log.Fatal(err)
	}

	r := pkg.NewRouter(c)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
