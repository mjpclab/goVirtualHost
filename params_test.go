package goVirtualHost

import (
	"crypto/tls"
	"testing"
)

func TestParamsValidateParam(t *testing.T) {
	var p *param
	var ps params
	var errs []error

	// normal wildcard ip
	p = &param{
		proto: "tcp",
		ip:    "",
		port:  ":80",
	}
	errs = ps.validateParam(p)
	if len(errs) > 0 {
		t.Error(errs)
	}
	ps = append(ps, p)

	// same wildcard ip:port, different hostname
	p = &param{
		proto:     "tcp",
		ip:        "",
		port:      ":80",
		hostNames: []string{"localhost"},
	}
	errs = ps.validateParam(p)
	if len(errs) > 0 {
		t.Error(errs)
	}
	ps = append(ps, p)

	// IPv4 wildcard 0.0.0.0:port, conflict
	p = &param{
		proto: "tcp",
		ip:    "0.0.0.0",
		port:  ":80",
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}

	// IPv6 wildcard [::]:port, conflict
	p = &param{
		proto: "tcp",
		port:  ":80",
	}
	p.ip = "[::]"
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}
	p.ip = "[::0]"
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}
	p.ip = "[0::0]"
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}
	p.ip = "[00::00]"
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}
	p.ip = "[0000:0000:0000:0000:0000:0000:0000:0000]"
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}

	// duplicated address and hostname
	p = &param{
		proto:     "tcp",
		ip:        "",
		port:      ":80",
		hostNames: []string{"localhost"},
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	}

	// cannot serve for both Plain and TLS mode
	p = &param{
		proto:  "tcp",
		ip:     "",
		port:   ":80",
		useTLS: true,
		cert:   &tls.Certificate{},
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error(errs)
	}
}
