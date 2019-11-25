package goVirtualHost

import "testing"

func TestExtractHostname(t *testing.T) {
	var host, hostname string

	host = "example.com"
	hostname = ExtractHostname(host)
	if hostname != "example.com" {
		t.Error(hostname)
	}

	host = "example.com:8080"
	hostname = ExtractHostname(host)
	if hostname != "example.com" {
		t.Error(hostname)
	}

	host = "[fe80::1]"
	hostname = ExtractHostname(host)
	if hostname != "[fe80::1]" {
		t.Error(hostname)
	}

	host = "[fe80::1]:8080"
	hostname = ExtractHostname(host)
	if hostname != "[fe80::1]" {
		t.Error(hostname)
	}
}

func TestNormalizeHostNames(t *testing.T) {
	inputs := []string{"aA", "", "Bb"}
	results := NormalizeHostNames(inputs)
	if len(results) != 2 || results[0] != "aa" || results[1] != "bb" {
		t.Error(results)
	}
}

func TestSplitListen(t *testing.T) {
	var proto string

	// ipv4
	ipv4 := "1.2.3.4"

	proto, ipv4Http := SplitListen(ipv4, false)
	if proto != "tcp4" {
		t.Error(proto)
	}
	if ipv4Http != "1.2.3.4:80" {
		t.Error(ipv4Http)
	}

	proto, ipv4Https := SplitListen(ipv4, true)
	if proto != "tcp4" {
		t.Error(proto)
	}
	if ipv4Https != "1.2.3.4:443" {
		t.Error(ipv4Https)
	}

	// ipv4:port
	ipv4Port := "2.3.4.5:6"

	proto, ipv4PortHttp := SplitListen(ipv4Port, false)
	if proto != "tcp4" {
		t.Error(proto)
	}
	if ipv4PortHttp != ipv4Port {
		t.Error(ipv4PortHttp)
	}

	proto, ipv4PortHttps := SplitListen(ipv4Port, true)
	if proto != "tcp4" {
		t.Error(proto)
	}
	if ipv4PortHttps != ipv4Port {
		t.Error(ipv4PortHttps)
	}

	// ipv6
	ipv6 := "[::1]"

	proto, ipv6Http := SplitListen(ipv6, false)
	if proto != "tcp6" {
		t.Error(proto)
	}
	if ipv6Http != "[::1]:80" {
		t.Error(ipv6Http)
	}

	proto, ipv6Https := SplitListen(ipv6, true)
	if proto != "tcp6" {
		t.Error(proto)
	}
	if ipv6Https != "[::1]:443" {
		t.Error(ipv6Https)
	}

	// ipv6:port
	ipv6Port := "[fe80::1234]:7"

	proto, ipv6PortHttp := SplitListen(ipv6Port, false)
	if proto != "tcp6" {
		t.Error(proto)
	}
	if ipv6PortHttp != ipv6Port {
		t.Error(ipv6PortHttp)
	}

	proto, ipv6PortHttps := SplitListen(ipv6Port, true)
	if proto != "tcp6" {
		t.Error(proto)
	}
	if ipv6PortHttps != ipv6Port {
		t.Error(ipv6PortHttps)
	}

	// port
	portNum := "8080"

	proto, portNumHttp := SplitListen(portNum, false)
	if proto != "tcp" {
		t.Error(proto)
	}
	if portNumHttp != ":8080" {
		t.Error(portNumHttp)
	}

	proto, portNumHttps := SplitListen(portNum, true)
	if proto != "tcp" {
		t.Error(proto)
	}
	if portNumHttps != ":8080" {
		t.Error(portNumHttps)
	}

	// :port
	port := ":3000"

	proto, portHttp := SplitListen(port, false)
	if proto != "tcp" {
		t.Error(proto)
	}
	if portHttp != port {
		t.Error(portHttp)
	}

	proto, portHttps := SplitListen(port, true)
	if proto != "tcp" {
		t.Error(proto)
	}
	if portHttps != port {
		t.Error(portHttps)
	}

	// hostname
	hostname := "example.com"

	proto, hostnameHttp := SplitListen(hostname, false)
	if proto != "tcp" {
		t.Error(proto)
	}
	if hostnameHttp != "example.com:80" {
		t.Error(hostnameHttp)
	}

	proto, hostnameHttps := SplitListen(hostname, true)
	if proto != "tcp" {
		t.Error(proto)
	}
	if hostnameHttps != "example.com:443" {
		t.Error(hostnameHttp)
	}

	// hostname:port
	host := "example.com:3210"

	proto, hostHttp := SplitListen(host, false)
	if proto != "tcp" {
		t.Error(proto)
	}
	if hostHttp != "example.com:3210" {
		t.Error(hostHttp)
	}

	proto, hostHttps := SplitListen(host, true)
	if proto != "tcp" {
		t.Error(proto)
	}
	if hostHttps != "example.com:3210" {
		t.Error(hostHttp)
	}
}
