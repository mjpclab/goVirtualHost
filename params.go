package goVirtualHost

import "errors"

func (params params) validate(param *param) error {
	if param.useTLS && param.cert == nil {
		return errors.New("certificate not found for TLS listens")
	}

	proto := param.proto
	addr := param.addr
	hostnames := param.hostNames

	for _, ownParam := range params {
		if ownParam.proto != proto || ownParam.addr != addr {
			continue
		}

		ownUseTLS := ownParam.cert != nil
		inputUseTLS := param.cert != nil
		if ownUseTLS != inputUseTLS {
			return errors.New("cannot served for both Plain and TLS mode")
		}

		if (len(hostnames) == 0 && len(ownParam.hostNames) == 0) || (ownParam.hasHostNames(hostnames)) {
			return errors.New("duplicated address and hostname")
		}
	}

	return nil
}
