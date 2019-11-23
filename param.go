package goVirtualHost

func (param *param) hasHostNames(checkHostNames []string) bool {
	if len(param.hostNames) == 0 || len(checkHostNames) == 0 {
		return false
	}

	for _, ownHostName := range param.hostNames {
		for _, checkHostName := range checkHostNames {
			if ownHostName == checkHostName {
				return true
			}
		}
	}
	return false
}
