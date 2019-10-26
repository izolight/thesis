package signingserver

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"net/url"
)

type Hash string
type IDPChoices struct {
	Urls              []url.URL
	InitialNonce      string
	IntermediateNonce string
}

const secret = "test1234"

var endpoint = url.URL{
	Host:   "idp.example.org",
	Scheme: "https",
	Path:   "authorize",
}
var redirectURL = url.URL{
	Scheme: "https",
	Host:   "localhost",
}

const clientID = "myClientID"

func PostHashes(req *restful.Request, resp *restful.Response) {
	var hashes []Hash
	err := req.ReadEntity(hashes)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	nonce, err := generateInitialNonce()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	intermediateNonce := generateIntermediateNonce(hashes, nonce, secret)
	OIDCNonce := generateOIDCNonce(hashes, intermediateNonce)
	endpointUrl := generateOIDCUrl(OIDCNonce, endpoint, clientID)
	choices := IDPChoices{
		Urls:              []url.URL{endpointUrl},
		InitialNonce:      nonce,
		IntermediateNonce: intermediateNonce,
	}
	resp.WriteEntity(choices)
}

func generateInitialNonce() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error reading random bytes: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

func generateIntermediateNonce(hashes []Hash, nonce string, secret string) string {
	hash := sha256.New()
	for _, h := range hashes {
		hash.Write([]byte(h))
	}
	hash.Write([]byte(nonce))
	hash.Write([]byte(secret))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func generateOIDCNonce(hashes []Hash, nonce string) string {
	hash := sha256.New()
	for _, h := range hashes {
		hash.Write([]byte(h))
	}
	hash.Write([]byte(nonce))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func generateOIDCUrl(nonce string, endpoint url.URL, clientId string) url.URL {
	q := endpoint.Query()
	q.Set("client_id", clientId)
	q.Add("response_type", "id_token")
	q.Add("scope", "openid")
	q.Add("nonce", nonce)
	q.Add("redirect_uri", redirectURL.String())
	endpoint.RawQuery = q.Encode()
	return endpoint
}