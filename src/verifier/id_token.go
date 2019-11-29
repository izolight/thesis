package verifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"gopkg.in/square/go-jose.v2"
	"time"
)

type idTokenVerifier struct {
	token []byte
	issuer string
	nonce string
	clientId string
	notAfter func() time.Time
	key jose.JSONWebKey
	ltv map[string]*LTV
	ctx context.Context
}

func NewIDTokenVerifier(token []byte, issuer, nonce, clientId string, notAfter time.Time, key []byte, ltv map[string]*LTV) (*idTokenVerifier, error) {
	i := &idTokenVerifier{
		token:    token,
		issuer:   issuer,
		nonce:    nonce,
		clientId: clientId,
		notAfter: notAfter.Local,
		ltv:      ltv,
		ctx: context.Background(),
		key: jose.JSONWebKey{},
	}
	err := i.key.UnmarshalJSON(key)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal jwk: %w", err)
	}
	return i, nil
}

func (i idTokenVerifier) VerifySignature(ctx context.Context, jwtRaw string) (payload []byte, err error) {
	signature, err := jose.ParseSigned(jwtRaw)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	return signature.Verify(i.key)
}

func (i idTokenVerifier) Verify() error {
	cfg := &oidc.Config{
		ClientID: i.clientId,
		Now: i.notAfter,
	}
	verifier := oidc.NewVerifier(i.issuer, i, cfg)
	idToken, err := verifier.Verify(i.ctx, string(i.token))
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
	l := ltvVerifier{
		certs:  i.key.Certificates,
		ltvMap: i.ltv,
	}
	err = l.Verify()
	if err != nil {
		return fmt.Errorf("ltv information for id token not valid: %w", err)
	}

	return nil
}

