package goVirtualHost

import (
	"./util"
	"crypto/tls"
	"net/http"
)

func newServer() *server {
	server := &server{
		vhosts:       vhosts{},
		defaultVhost: nil,

		httpServer: &http.Server{},
	}

	return server
}

func (server *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var vhost *vhost

	hostname := util.ExtractHostname(r.Host)

	for _, vh := range server.vhosts {
		if vh.matchHostName(hostname) {
			vhost = vh
			break
		}
	}

	if vhost == nil {
		vhost = server.defaultVhost
	}

	vhost.handler.ServeHTTP(w, r)
}

func (server *server) updateDefaultVhost() {
	for _, vh := range server.vhosts {
		if len(vh.hostNames) == 0 {
			server.defaultVhost = vh
			break
		}
	}

	if server.defaultVhost == nil {
		server.defaultVhost = server.vhosts[0]
	}
}

func (server *server) updateHttpServerTLSConfig() {
	certs := []tls.Certificate{}
	for _, vhost := range server.vhosts {
		if vhost.cert == nil {
			continue
		}
		certs = append(certs, *vhost.cert)
	}

	if len(certs) == 0 {
		server.httpServer.TLSConfig = nil
		return
	}

	tlsConfig := &tls.Config{
		Certificates: certs,
	}
	tlsConfig.BuildNameToCertificate()
	server.httpServer.TLSConfig = tlsConfig
}

func (server *server) updateHttpServerHandler() {
	if len(server.vhosts) == 1 {
		server.httpServer.Handler = server.defaultVhost.handler
		return
	}

	server.httpServer.Handler = server
}

func (server *server) open(listener *listener) error {
	if server.httpServer.TLSConfig != nil {
		return server.httpServer.ServeTLS(listener.netListener, "", "")
	} else {
		return server.httpServer.Serve(listener.netListener)
	}
}

func (server *server) close() error {
	return server.httpServer.Close()
}
