package sscert

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	rand2 "math/rand"
	"micro/util"
	"net"
	"os"
	"strings"
	"time"
)

func EnsureCertificatesCreated(hostname string) {
	if util.FileExists("server.crt") && util.FileExists("server.key") {
		return
	} else {
		NewSSCertificate(hostname)
	}
}

func NewSSCertificate(hostname string) {
	log.Println("Creating new self signed certificate")
	//privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	util.LogFatal(err)
	template := x509.Certificate{
		SerialNumber: big.NewInt(rand2.Int63()),
		Subject: pkix.Name{
			CommonName:         "web",
			OrganizationalUnit: []string{"IT"},
			Organization:       []string{"Prasenjit.Net"},
			Country:            []string{"IN"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(hostname, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey)
	util.LogFatalWithMessage("Failed to create certificate: %s", err)

	var certFile *os.File
	var keyFile *os.File
	if util.FileExists("server.crt") {
		certFile, err = os.Open("server.crt")
	} else {
		certFile, err = os.Create("server.crt")
	}
	util.LogFatalWithMessage("failed to open cert file", err)
	defer util.Close(certFile)
	if util.FileExists("server.key") {
		keyFile, err = os.Open("server.key")
	} else {
		keyFile, err = os.Create("server.key")
	}
	util.LogFatalWithMessage("failed to open key file", err)
	defer util.Close(keyFile)

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	util.LogFatalWithMessage("failed to write cert file", err)
	err = pem.Encode(keyFile, pemBlockForKey(privateKey))
	util.LogFatalWithMessage("failed to write key file", err)

}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(key interface{}) *pem.Block {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		util.LogFatalWithMessage("Unable to marshal ECDSA private key: %v", err)
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
