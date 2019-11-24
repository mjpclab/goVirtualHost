package goVirtualHost

import (
	"./util"
	"crypto/tls"
	"strings"
)

func (info *HostInfo) toParam(listen string, useTLS bool) *param {
	proto, addr := util.SplitListen(listen, false)
	var cert *tls.Certificate
	if useTLS {
		cert = info.Cert
	}

	param := &param{
		proto:   proto,
		addr:    addr,
		useTLS:  useTLS,
		cert:    cert,
		handler: info.Handler,
	}

	return param
}

func (info *HostInfo) toParams() params {
	params := params{}

	hostNames := util.FilterEmptyStrings(info.HostNames)
	for i, s := 0, len(hostNames); i < s; i++ {
		hostNames[i] = strings.ToLower(hostNames[i])
	}

	for _, listen := range info.Listens {
		param := info.toParam(listen, info.Cert != nil)
		param.hostNames = hostNames
		params = append(params, param)
	}

	for _, listen := range info.ListensPlain {
		param := info.toParam(listen, false)
		param.hostNames = hostNames
		params = append(params, param)
	}

	for _, listen := range info.ListensTLS {
		param := info.toParam(listen, true)
		param.hostNames = hostNames
		params = append(params, param)
	}

	return params
}
