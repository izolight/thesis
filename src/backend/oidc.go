package backend

import (
	"context"
	"golang.org/x/oauth2"
	oidc "github.com/coreos/go-oidc"
	"net/http"
)

type OIDCClient struct {
	ctx context.Context
	*oidc.Provider
	*oauth2.Config
}

func newOIDCProvider(providerURL, clientID, clientSecret, redirectURL string) (*OIDCClient, error) {
	client := &OIDCClient{
		ctx: context.Background(), // TODO: change to real context
	}
	provider, err := oidc.NewProvider(client.ctx, providerURL)
	if err != nil {
		return nil, err
	}
	client.Provider = provider

	config := &oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		RedirectURL: redirectURL,
		Endpoint: client.Provider.Endpoint(),
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	client.Config = config

	return client, nil
}

func (c *OIDCClient) oidcRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.Config.AuthCodeURL(state), http.StatusFound)
}

func (c *OIDCClient) oidcCallback(w http.ResponseWriter, r *http.Request) {
	oauth2Token, err := c.Config.Exchange(c.ctx, r.URL.Query().Get("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse id token"))
		return
	}

	idToken, err := verifier.Verify(c.ctx, rawIDToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}