package main

import (
	goVirtualHost "../.."
	"fmt"
	"net/http"
)

func main() {
	println("Starting...")

	var errors []error
	var err error
	svc := goVirtualHost.NewService()

	// virtual host: localhost
	certLocalhost, err := goVirtualHost.LoadCertificate("example/tls/localhost.crt", "example/tls/localhost.key")
	if err != nil {
		fmt.Println(err)
	}
	errors = svc.Add(&goVirtualHost.HostInfo{
		Listens:   []string{":8080"},
		Cert:      certLocalhost,
		HostNames: []string{"localhost"},
		Handler:   http.FileServer(http.Dir(".")),
	})
	if len(errors) > 0 {
		fmt.Println(errors)
	}

	// virtual host: default
	certExample, err := goVirtualHost.LoadCertificate("example/tls/example.crt", "example/tls/example.key")
	if err != nil {
		fmt.Println(err)
	}
	errors = svc.Add(&goVirtualHost.HostInfo{
		Listens:   []string{":8080"},
		Cert:      certExample,
		HostNames: nil,
		Handler:   http.FileServer(http.Dir("/tmp")),
	})
	if len(errors) > 0 {
		fmt.Println(errors)
	}

	// start server
	errors = svc.Open()
	if len(errors) > 0 {
		fmt.Println(errors)
	}
}
