package verifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"time"
)

type idTokenVerifier struct {
	token []byte
	issuer string
	keys string
	nonce string
	clientId string
	notAfter time.Time
}

func (i idTokenVerifier) Verify() error {
	ctx := context.Background()
	keySet := oidc.NewRemoteKeySet(ctx, i.keys)
	cfg := &oidc.Config{
		ClientID: i.clientId,
		Now: i.notAfter.Local,
	}
	verifier := oidc.NewVerifier(i.issuer, keySet, cfg)
	idToken, err := verifier.Verify(ctx, string(i.token))
	if err != nil {
		return err
	}
	if idToken.Nonce != i.nonce {
		return fmt.Errorf("nonce didn't match, was %s, should be :%s", idToken.Nonce, i.nonce)
	}
	var emailClaims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err = idToken.Claims(&emailClaims); err != nil {
		return err
	}
	if !emailClaims.EmailVerified {
		return errors.New("e-mail was not verified")
	}

	return nil
}