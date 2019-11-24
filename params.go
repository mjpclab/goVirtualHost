package goVirtualHost

import "fmt"

func (params params) validate(param *param) error {
	if param.useTLS && param.cert == nil {
		return fmt.Errorf("certificate not found for TLS listens: %+v", param)
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
			return fmt.Errorf("cannot serve for both Plain and TLS mode: %+v", param)
		}

		if (len(hostnames) == 0 && len(ownParam.hostNames) == 0) || (ownParam.hasHostNames(hostnames)) {
			return fmt.Errorf("duplicated address and hostname: %+v", param)
		}
	}

	return nil
}
