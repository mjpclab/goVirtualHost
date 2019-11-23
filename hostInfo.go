package goVirtualHost

import (
	"./util"
	"strings"
)

func (info *HostInfo) toParams() params {
	params := params{}

	hostNames := util.FilterEmptyStrings(info.HostNames)
	for i, s := 0, len(hostNames); i < s; i++ {
		hostNames[i] = strings.ToLower(hostNames[i])
	}

	for _, listen := range info.Listens {
		proto, addr := util.SplitListen(listen, info.Cert != nil)

		param := &param{
			proto:     proto,
			addr:      addr,
			cert:      info.Cert,
			hostNames: hostNames,
			handler:   info.Handler,
		}
		params = append(params, param)
	}

	return params
}
