package goVirtualHost

import (
	"crypto/tls"
	"net"
	"net/http"
)

// init host info
type HostInfo struct {
	Listens   []string
	Cert      *tls.Certificate
	HostNames []string
	Handler   http.Handler
}

// normalized HostInfo Param
type param struct {
	proto     string
	addr      string
	cert      *tls.Certificate
	hostNames []string
	handler   http.Handler
}

type params []*param

// wrapper of net.Listener
type listener struct {
	proto       string
	addr        string
	netListener net.Listener
	server      *server
}

type listeners []*listener

// wrapper for http.Server
type server struct {
	vhosts       vhosts
	defaultVhost *vhost

	dispatchHandler http.Handler
	httpServer      *http.Server
}

type servers []*server

// virtual host
type vhost struct {
	cert      *tls.Certificate
	hostNames []string
	handler   http.Handler
}

type vhosts []*vhost

// service
type Service struct {
	params    params
	listeners listeners
	servers   servers
	vhosts    vhosts
}
