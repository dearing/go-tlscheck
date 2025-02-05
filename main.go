package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var version = "1.1.2"

var argUrl = flag.String("url", "https://www.google.com", "URL to connect to")
var argJson = flag.Bool("json", false, "Print output as JSON")
var argVersion = flag.Bool("version", false, "Print version")

func main() {

	flag.Parse()

	if *argVersion {
		fmt.Printf("github.com/dearing/go-tlscheck v%s\n", version)
		return
	}

	startTime := time.Now()

	req, err := http.Get(*argUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	if *argJson {
		jsonPrint(req.TLS)
		return
	}

	fmt.Printf("GET %s took %s\n", *argUrl, time.Since(startTime))

	subject := req.TLS.PeerCertificates[0].Subject

	fmt.Printf("CommonName:    %s\n", subject.CommonName)
	fmt.Printf("DNSNames:      %s\n", req.TLS.PeerCertificates[0].DNSNames)
	fmt.Printf("IssuerOrg:     %s\n", req.TLS.PeerCertificates[0].Issuer.Organization)

	for _, names := range subject.Names {
		fmt.Printf("Subject:       %s\n", names.Value)
	}

	cipherSuite := req.TLS.CipherSuite
	fmt.Printf("CipherSuite:   %s\n", tls.CipherSuiteName(cipherSuite))

	fmt.Printf("NotBefore:     %s\n", req.TLS.PeerCertificates[0].NotBefore)
	fmt.Printf("NotAfter:      %s\n", req.TLS.PeerCertificates[0].NotAfter)
	
	expiresIn := req.TLS.PeerCertificates[0].NotAfter.Sub(time.Now())
	fmt.Printf("ExpiresIn:     %s\n", expiresIn)
}

func jsonPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
