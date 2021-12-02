package goVirtualHost

import (
	"errors"
	"net"
)

var unknownIPVersion = errors.New("unknown IP version")

func newIPAddr(netIP net.IP) (*ipAddr, error) {
	var version int
	if netIP.To4() != nil {
		version = ip4ver
	} else if netIP.To16() != nil {
		version = ip6ver
	} else {
		return nil, unknownIPVersion
	}

	instance := &ipAddr{
		netIP:              netIP,
		version:            version,
		isGlobalUnicast:    netIP.IsGlobalUnicast(),
		isLinkLocalUnicast: netIP.IsLinkLocalUnicast(),
	}
	return instance, nil
}

func (addr *ipAddr) String() string {
	if addr.version == ip6ver {
		return "[" + addr.netIP.String() + "]"
	} else {
		return addr.netIP.String()
	}
}
