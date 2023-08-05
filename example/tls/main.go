package main

import (
	"crypto/tls"
	"fmt"
	"mjpclab.dev/goVirtualHost"
	"net/http"
)

func main() {
	println("Starting...")

	var errors []error
	var err error
	svc := goVirtualHost.NewService()

	// virtual host: localhost
	certLocalhost, err := goVirtualHost.LoadCertificate("localhost.crt", "localhost.key")
	if err != nil {
		fmt.Println(err)
	}
	errors = svc.Add(&goVirtualHost.HostInfo{
		Listens:   []string{":8080"},
		Certs:     []tls.Certificate{certLocalhost},
		HostNames: []string{"localhost"},
		Handler:   http.FileServer(http.Dir(".")),
	})
	if len(errors) > 0 {
		fmt.Println(errors)
	}

	// virtual host: default
	certExample, err := goVirtualHost.LoadCertificate("example.crt", "example.key")
	if err != nil {
		fmt.Println(err)
	}
	errors = svc.Add(&goVirtualHost.HostInfo{
		Listens:   []string{":8080"},
		Certs:     []tls.Certificate{certExample},
		HostNames: nil,
		Handler:   http.FileServer(http.Dir("/tmp")),
	})
	if len(errors) > 0 {
		fmt.Println(errors)
	}

	// print accessible urls
	vhostsUrls := svc.GetAccessibleURLs(false)
	fmt.Println("Accessible urls:")
	for vhIndex, urls := range vhostsUrls {
		fmt.Println("virtual host", vhIndex)
		for _, url := range urls {
			fmt.Println("  ", url)
		}
	}

	// start server
	errors = svc.Open()
	if len(errors) > 0 {
		fmt.Println(errors)
	}
}
