package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var version = "1.1.3"

var argUrl = flag.String("url", "https://www.google.com", "URL to connect to")
var argTimeout = flag.Int("timeout", 300, "timeout in seconds")
var argJson = flag.Bool("json", false, "write to STDOUT certificate information JSON")
var argVersion = flag.Bool("version", false, "output version and return")

func usage() {
	println(`Usage: go-tlscheck [options]
Quickly check the TLS certificate of a website and display useful information
or enable JSON output to get all the information in from the Go handler.

- ex: go-tlscheck -url https://www.google.com
- ex: go-tlscheck -url https://www.google.com -json | jq '.PeerCertificates[].SerialNumber'

Options:
`)
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *argVersion {
		println("go-tlscheck v" + version)
		info, ok := debug.ReadBuildInfo()
		if ok {
			slog.Info("build info", "main", info.Main.Path, "version", info.Main.Version)
			for _, setting := range info.Settings {
				slog.Info("build info", "key", setting.Key, "value", setting.Value)
			}
		}

		fmt.Printf("github.com/dearing/go-tlscheck v%s\n", version)
		return
	}

	startTime := time.Now()

	httpClient := &http.Client{
		Timeout: time.Duration(*argTimeout) * time.Second,
	}

	req, err := httpClient.Get(*argUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	if *argJson {
		jsonPrint(req.TLS)
		return
	}

	slog.Info("GET", "url", *argUrl, "duration", time.Since(startTime))

	if err != nil {
		slog.Error("GET", "url", *argUrl, "error", err)
		os.Exit(1)
	}
	defer req.Body.Close()

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

	expiresIn := time.Until(req.TLS.PeerCertificates[0].NotAfter)
	fmt.Printf("ExpiresIn:     %s\n", expiresIn)
}

func jsonPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
