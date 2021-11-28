package goVirtualHost

import (
	"crypto/tls"
	"testing"
)

func TestParamValidate(t *testing.T) {
	var p *param
	var errs []error

	p = &param{
		proto: "tcp",
		ip:    "",
		port:  "80",
	}
	errs = p.validate()
	if len(errs) > 0 {
		t.Error()
	}

	p.useTLS = true
	errs = p.validate()
	if len(errs) == 0 {
		t.Error()
	}

	p.cert = &tls.Certificate{}
	errs = p.validate()
	if len(errs) > 0 {
		t.Error()
	}
}
