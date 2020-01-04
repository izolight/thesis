package verifier

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/sirupsen/logrus"
)


type Config struct {
	Issuer          string
	ClientId        string
	AdditionalCerts []*x509.Certificate
	Logger          *logrus.Entry
}

func NewDefaultCfg(caFile []byte) Config {
	cfg := Config{
		Issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
		ClientId: "thesis",
	}
	if caFile == nil {
		return cfg
	}
	filePEM, _ := pem.Decode(caFile)
	if filePEM == nil {
		logrus.Println("couldn't decode pem")
		return cfg
	}
	rootCA, err := x509.ParseCertificate(filePEM.Bytes)
	if err != nil {
		logrus.Println(err)
		return cfg
	}
	cfg.AdditionalCerts = []*x509.Certificate{
		rootCA,
	}
	return cfg
}