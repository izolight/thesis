package main

import (
	"github.com/sirupsen/logrus"
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

	logger := logrus.New()

	r := backend.NewRouter(c, logger)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Infof("Started server at %s", srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
