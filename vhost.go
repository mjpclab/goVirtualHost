package goVirtualHost

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"strings"
)

func makeLeaf(certs certs) {
	for _, cert := range certs {
		if cert.Leaf == nil && len(cert.Certificate) > 0 {
			if leaf, err := x509.ParseCertificate(cert.Certificate[0]); err == nil {
				cert.Leaf = leaf
			}
		}
	}
}

func newVhost(hostNames []string, certKeyPaths certKeyPairs, vhCerts certs, handler http.Handler) *vhost {
	makeLeaf(vhCerts)

	vhost := &vhost{
		hostNames:    hostNames,
		certKeyPaths: certKeyPaths,
		certs:        vhCerts,
		loadedCerts:  vhCerts,
		handler:      handler,
	}

	return vhost
}

func (vh *vhost) matchHostName(name string) bool {
	reqHostName := strings.ToLower(name)
	for _, hostname := range vh.hostNames {
		if hostname == reqHostName {
			return true
		}
		if len(hostname) > 1 {
			if hostname[0] == '.' && strings.HasSuffix(reqHostName, hostname) {
				return true
			} else if hostname[len(hostname)-1] == '.' && strings.HasPrefix(reqHostName, hostname) {
				return true
			}
		}
	}
	return false
}

func (vh *vhost) loadCertificates() []error {
	loadedCerts, errs := LoadCertificatesFromPairs(vh.certKeyPaths)
	makeLeaf(loadedCerts)
	loadedCerts = append(loadedCerts, vh.certs...)

	vh.loadedCerts = loadedCerts

	return errs
}

func (vh *vhost) lookupCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certs := vh.loadedCerts
	certLen := len(certs)
	if certLen == 1 {
		return certs[0], nil
	}

	for _, cert := range certs {
		if cert.Leaf == nil {
			continue
		}
		err := cert.Leaf.VerifyHostname(hello.ServerName)
		if err == nil {
			return cert, err
		}
	}

	if certLen > 0 {
		return certs[0], nil
	}

	return nil, errors.New("cannot find proper certificate for " + hello.ServerName)
}
