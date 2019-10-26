package signingserver

import (
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"github.com/coreos/go-oidc"
	"github.com/emicklei/go-restful"
	"golang.org/x/oauth2"
	"net/http"
)

type SigningRequest struct {
	JWT oauth2.Token
	//IDToken oidc.IDToken
	Hashes []Hash
	InitialNonce string
	IntermediateNonce string
}

type IntermediateSigningResponse struct {
	SignedHashes [][]byte
	IDToken string
	IntermediateNonce string
	Certificate []byte
}

type SigningResponse struct {
	IntermediateSigningResponse
	Timestamp []byte
}

func SignHashes(req *restful.Request, resp *restful.Response) {
	signReq := new(SigningRequest)
	err := req.ReadEntity(&signReq)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	intermediateNonce := generateIntermediateNonce(signReq.Hashes, signReq.InitialNonce, secret)
	if intermediateNonce != signReq.IntermediateNonce {
		resp.WriteErrorString(http.StatusBadRequest, "Intermediate nonce has been manipulated")
		return
	}
	rawIDToken, ok := signReq.JWT.Extra("id_token").(string)
	if !ok {
		resp.WriteErrorString(http.StatusBadRequest, "id_token missing")
		return
	}
	idToken, err := verifier.Verify(req.Request.Context(), rawIDToken)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, "Failed to verify ID Token: "+err.Error())
		return
	}

	OIDCNonce := generateOIDCNonce(signReq.Hashes, intermediateNonce)
	if OIDCNonce != idToken.Nonce {
		resp.WriteErrorString(http.StatusBadRequest, "OIDC nonce has been manipulated")
		return
	}

	csr, err := requestSigningKeyCSR(idToken)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	cert, err := requestSigningKeySigning(csr, true)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	signatures, err := requestSignatures(signReq.Hashes, idToken.Subject)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	intermediateResponse := IntermediateSigningResponse{
		SignedHashes:      signatures,
		IDToken:           rawIDToken,
		IntermediateNonce: signReq.IntermediateNonce,
		Certificate: cert,
	}

	timestamp, err := requestTimestamp(intermediateResponse)

	response := SigningResponse{
		intermediateResponse,
		timestamp,
	}
	resp.WriteEntity(response)
}

func requestSigningKeyCSR(token *oidc.IDToken) ([]byte, error) {
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: token.Subject,
			Country: []string{"CH"},
			Province: []string{"BE"},
			Locality: []string{"Bern"},
			Organization: []string{"SigningService"},
			OrganizationalUnit: []string{"Testing"},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	return generateLocalKey(template)
}

func requestSigningKeySigning(csr []byte, signLocal bool) ([]byte, error) {
	if signLocal {
		return signCSRWithLocalCA(csr)
	}
	return nil, errors.New("Only local signing is implemented for now")
}

func requestSignatures(hashes []Hash, subject string) ([][]byte, error) {
	signatures := [][]byte{}
	for _, hash := range hashes {
		signature, err := sign(hash, subject)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, signature)
	}
	return signatures, nil
}

func requestTimestamp(response IntermediateSigningResponse) ([]byte, error) {
	hash := sha256.New()
	for _, sig := range response.SignedHashes {
		hash.Write(sig)
	}
	hash.Write([]byte(response.IntermediateNonce))
	hash.Write([]byte(response.IDToken))
	// TODO:
	return hash.Sum(nil), nil
}