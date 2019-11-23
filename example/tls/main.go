package main

import (
	goVirtualHost "../.."
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	println("Starting...")

	var err error
	svc := goVirtualHost.NewService()

	// virtual host: localhost
	certLocalhost, err := tls.LoadX509KeyPair("example/tls/localhost.crt", "example/tls/localhost.key")
	if err != nil {
		fmt.Println(err)
	}
	err = svc.Add(&goVirtualHost.HostInfo{
		Listen:    ":8080",
		Cert:      &certLocalhost,
		HostNames: []string{"localhost"},
		Handler:   http.FileServer(http.Dir(".")),
	})
	if err != nil {
		fmt.Println(err)
	}

	// virtual host: default
	certExample, err := tls.LoadX509KeyPair("example/tls/example.crt", "example/tls/example.key")
	if err != nil {
		fmt.Println(err)
	}
	err = svc.Add(&goVirtualHost.HostInfo{
		Listen:    ":8080",
		Cert:      &certExample,
		HostNames: nil,
		Handler:   http.FileServer(http.Dir("/tmp")),
	})
	if err != nil {
		fmt.Println(err)
	}

	// start server
	svc.Open()
}
