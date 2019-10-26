package signingserver

import (
	"context"
	"github.com/coreos/go-oidc"
	"log"
)

var provider oidc.Provider
var verifier *oidc.IDTokenVerifier

func OIDCInit() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier = provider.Verifier(oidcConfig)
}