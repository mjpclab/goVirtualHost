package goVirtualHost

import (
	"./util"
	"strings"
)

func (info *HostInfo) toParam() *param {
	proto, addr := util.SplitListen(info.Listen, info.Cert != nil)

	hostNames := util.FilterEmptyStrings(info.HostNames)
	for i, s := 0, len(hostNames); i < s; i++ {
		hostNames[i] = strings.ToLower(hostNames[i])
	}

	param := &param{
		proto:     proto,
		addr:      addr,
		cert:      info.Cert,
		hostNames: hostNames,
		handler:   info.Handler,
	}

	return param
}
