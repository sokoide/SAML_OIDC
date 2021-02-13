//main.go
package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", samlsp.AttributeFromContext(r.Context(), "username"))
}

func main() {
	samlIdpMetadataURL := "https://timemachine:8443/auth/realms/master/protocol/saml/descriptor"
	confRootURL := "http://scottmm.local:8000"
	confPort := ":8000"

	keyPair, err := tls.LoadX509KeyPair("service.crt", "service.key")
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}

	// idpMetadataURL, err := url.Parse("https://samltest.id/saml/idp")
	idpMetadataURL, err := url.Parse(samlIdpMetadataURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	// Allow self-signed cert
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	// idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient,
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), httpClient,
		*idpMetadataURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	rootURL, err := url.Parse(confRootURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})
	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	http.ListenAndServe(confPort, nil)

}
