package main

import (
	"fmt"
	"github.com/izolight/go_S-MIME/timestamp"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide a file")
	}
	in := os.Args[1]
	data, err := ioutil.ReadFile(in)
	if err != nil {
		log.Fatalf("Could not read file %s: %s", in, err)
	}

	tsr, err := timestamp.ParseResponse(data)
	if err != nil {
		log.Fatalf("Could not parse timestamp: %s", err)
	}

	tsinfo, err := tsr.Info()
	if err != nil {
		log.Fatalf("Could not get tsinfo: %s", err)
	}

	contentInfo := tsr.TimeStampToken
	//fmt.Println(contentInfo.ContentType)

	signedData, err := contentInfo.SignedDataContent()
	if err != nil {
		log.Fatalf("could not get signed data: %s", err)
	}

	certs, err := signedData.X509Certificates()
	if err != nil {
		log.Fatalf("could not get x509 certs: %s", err)
	}

	for _, s := range signedData.SignerInfos {
		fmt.Println(s.X509SignatureAlgorithm())
	}
	fmt.Println(signedData.EncapContentInfo.EContentType)

	fmt.Println("TimeStampResp")
	fmt.Printf("\tStatus: %s\n", tsr.Status)
	fmt.Println("\tTimeStampToken:")
	fmt.Printf("\t\tcontentType: %s\n", contentInfo.ContentType)
	fmt.Println("\t\tcontent:")
	fmt.Printf("\t\t\tVersion: %d\n", signedData.Version)
	fmt.Println("\t\t\tDigestAlgorithms:")
	for _, dA := range signedData.DigestAlgorithms {
		fmt.Printf("\t\t\t\tAlgorithm: %s\n", dA.Algorithm)
	}
	fmt.Println("\t\t\tEncapContentInfo:")
	fmt.Printf("\t\t\t\tEContentType: %s\n", signedData.EncapContentInfo.EContentType)
	fmt.Println("\t\t\t\tEContent:")
	fmt.Printf("\t\t\t\t\tVersion: %d\n", tsinfo.Version)
	fmt.Printf("\t\t\t\t\tPolicy: %s\n", tsinfo.Policy)
	fmt.Println("\t\t\t\t\tMessageImprint:")
	fmt.Printf("\t\t\t\t\t\tHashAlgorithm: %s\n", tsinfo.MessageImprint.HashAlgorithm.Algorithm)
	fmt.Printf("\t\t\t\t\t\tHashedMessage: %x\n", tsinfo.MessageImprint.HashedMessage)
	fmt.Printf("\t\t\t\t\tSerialNumber: %s\n", tsinfo.SerialNumber)
	fmt.Printf("\t\t\t\t\tGenTime: %s\n", tsinfo.GenTime)
	fmt.Println("\t\t\t\t\tAccuracy:")
	fmt.Printf("\t\t\t\t\t\tSeconds: %d\n", tsinfo.Accuracy.Seconds)
	fmt.Printf("\t\t\t\t\t\tMillis: %d\n", tsinfo.Accuracy.Millis)
	fmt.Printf("\t\t\t\t\t\tMicros: %d\n", tsinfo.Accuracy.Micros)
	fmt.Printf("\t\t\t\t\tOrdering: %t\n", tsinfo.Ordering)
	fmt.Printf("\t\t\t\t\tNonce: %s\n", tsinfo.Nonce)
	fmt.Println("\t\t\t\t\tTSA:")
	fmt.Printf("\t\t\t\t\t\t%v\n", tsinfo.TSA)
	fmt.Println("\t\t\t\t\tExtensions:")
	for _, ex := range tsinfo.Extensions {
		fmt.Printf("\t\t\t\t\t\tId: %s\n", ex.Id)
		fmt.Printf("\t\t\t\t\t\tCritical: %t\n", ex.Critical)
		fmt.Printf("\t\t\t\t\t\tValue: %x\n", ex.Value)
	}
	fmt.Println("\t\t\tCertificates:")
	for _, cert := range certs {
		fmt.Printf("\t\t\t\tSubject: %s\n", cert.Subject)
		fmt.Printf("\t\t\t\tIssuer: %s\n", cert.Issuer)
		fmt.Printf("\t\t\t\tIssuerURL: %s\n", cert.IssuingCertificateURL)
	}
	fmt.Printf("\t\t\tCRLs: %v\n", signedData.CRLs)
	fmt.Println("\t\t\tSignerInfos:")
	for _, si := range signedData.SignerInfos {
		fmt.Printf("\t\t\t\tVersion: %d\n", si.Version)
		fmt.Printf("\t\t\t\tDigestAlgorithm: %s\n", si.DigestAlgorithm.Algorithm)
		fmt.Printf("\t\t\t\tSignatureAlgorithm: %s\n", si.SignatureAlgorithm.Algorithm)
		fmt.Println("\t\t\t\tSID:")
		fmt.Printf("\t\t\t\t\tSerialNumber: %s\n", si.SID.IAS.SerialNumber)
		fmt.Printf("\t\t\t\t\tIssuer: %v\n", si.SID.IAS.Issuer)
		fmt.Printf("\t\t\t\tSignature: %x\n", si.Signature)
		fmt.Println("\t\t\t\tSignedAttributes:")
		for _, sa := range si.SignedAttrs {
			fmt.Printf("\t\t\t\t\tType:%s\n", sa.Type)
		}
		fmt.Println("\t\t\t\tUnsignedAttributes:")
		for _, ua := range si.UnsignedAttrs {
			fmt.Printf("\t\t\t\t\tType:%s\n", ua.Type)
		}
		h, _ := si.Hash()
		fmt.Printf("\t\t\t\tHash: %v\n", h)
		ct, _ := si.GetContentTypeAttribute()
		fmt.Printf("\t\t\t\tContent-Type: %s\n", ct)
		dig, _ := si.GetMessageDigestAttribute()
		fmt.Printf("\t\t\t\tDigest: %x\n", dig)
		st, _ := si.GetSigningTimeAttribute()
		fmt.Printf("\t\t\t\tTime: %s\n", st)
	}
	_, err = timestamp.VerfiyTS(contentInfo)
	if err != nil {
		log.Fatal(err)
	}
}
