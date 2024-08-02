package goVirtualHost

import (
	"net/http"
	"strings"
)

func newVhost(hostNames []string, certKeyPaths certKeyPairs, certs certs, handler http.Handler) *vhost {
	vhost := &vhost{
		hostNames:    hostNames,
		certKeyPaths: certKeyPaths,
		certs:        certs,
		handler:      handler,
	}

	return vhost
}

func (v *vhost) matchHostName(name string) bool {
	reqHostName := strings.ToLower(name)
	for _, hostname := range v.hostNames {
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
