# go-tlscheck

quick and dirty TLS cert information

```shell
go install github.com/dearing/go-tlscheck@latest
```
---
```shell
go-tlscheck -h
Usage of go-tlscheck:
  -json
        Print output as JSON
  -url string
        URL to connect to (default "https://www.google.com")
  -version
        Print version
```
---
```
go-tlscheck -url https://github.com
CommonName:    github.com
DNSNames:      [github.com www.github.com]
IssuerOrg:     [Sectigo Limited]
Subject:       github.com
CipherSuite:   TLS_AES_128_GCM_SHA256
NotBefore:     2024-03-07 00:00:00 +0000 UTC
NotAfter:      2025-03-07 23:59:59 +0000 UTC
ExpiresIn:     839h29m17.232334432s

```

