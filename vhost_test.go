package goVirtualHost

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sync"
	"testing"
	"time"
)

func TestMatchHostName(t *testing.T) {

	var vh *vhost

	vh = newVhost([]string{"www.example.com"}, nil, nil, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost([]string{".example.com"}, nil, nil, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost([]string{".example.com", "example.com"}, nil, nil, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if !vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost([]string{"example."}, nil, nil, nil)
	if !vh.matchHostName("example.com") {
		t.Error()
	}
	if !vh.matchHostName("example.net") {
		t.Error()
	}
}

func newTestCert(t *testing.T, dnsName string) *tls.Certificate {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: dnsName},
		DNSNames:     []string{dnsName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return &tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

func TestLookupCertificate(t *testing.T) {
	cert1 := newTestCert(t, "a.example.com")
	cert2 := newTestCert(t, "b.example.com")
	cert3 := newTestCert(t, "c.example.com")
	cert3.Leaf, _ = x509.ParseCertificate(cert3.Certificate[0])
	vh := newVhost(nil, nil, certs{cert1, cert2, cert3}, nil)
	vh.loadCertificates()

	if cert, _ := vh.lookupCertificate(&tls.ClientHelloInfo{ServerName: "c.example.com"}); cert != cert3 {
		t.Error("certificate with Leaf should be kept as is")
	}

	for _, serverName := range []string{"a.example.com", "b.example.com"} {
		cert, err := vh.lookupCertificate(&tls.ClientHelloInfo{ServerName: serverName})
		if err != nil {
			t.Fatal(err)
		}
		if cert.Leaf == nil || cert.Leaf.Subject.CommonName != serverName {
			t.Error(serverName, cert.Leaf)
		}
	}
}

func TestLookupCertificateConcurrently(t *testing.T) {
	vh := newVhost(nil, nil, certs{newTestCert(t, "a.example.com"), newTestCert(t, "b.example.com")}, nil)
	vh.loadCertificates()

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cert, err := vh.lookupCertificate(&tls.ClientHelloInfo{ServerName: "b.example.com"})
			if err != nil || cert == nil {
				t.Error(cert, err)
			}
		}()
	}
	wg.Wait()
}
