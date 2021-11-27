package goVirtualHost

import "fmt"

func (params params) validateParam(param *param) (errs []error) {
	for _, ownParam := range params {
		if ownParam == param {
			continue
		}

		if ownParam.port == param.port && ownParam.proto != "unix" && param.proto != "unix" {
			if (ownParam.proto == "tcp" && ownParam.ip == "") ||
				(param.proto == "tcp" && param.ip == "") ||
				(ownParam.proto == param.proto && (ownParam.ip == "" || param.ip == "")) {
				err := fmt.Errorf("conflict IP address: %+v, %+v", ownParam, param)
				errs = append(errs, err)
			}
		}

		if ownParam.proto == param.proto && ownParam.ip == param.ip && ownParam.port == param.port {
			ownUseTLS := ownParam.cert != nil
			useTLS := param.cert != nil
			if ownUseTLS != useTLS {
				err := fmt.Errorf("cannot serve for both Plain and TLS mode: %+v", param)
				errs = append(errs, err)
			}

			if (len(param.hostNames) == 0 && len(ownParam.hostNames) == 0) || (ownParam.hasHostNames(param.hostNames)) {
				err := fmt.Errorf("duplicated address and hostname: %+v", param)
				errs = append(errs, err)
			}
		}
	}

	return
}

func (params params) validate(inputs params) (errs []error) {
	for _, p := range inputs {
		es := p.validate()
		if len(es) > 0 {
			errs = append(errs, es...)
		}

		es = inputs.validateParam(p)
		if len(es) > 0 {
			errs = append(errs, es...)
		}

		es = params.validateParam(p)
		if len(es) > 0 {
			errs = append(errs, es...)
		}
	}

	return
}
