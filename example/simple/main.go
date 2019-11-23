package main

import (
	goVirtualHost "../.."
	"fmt"
	"net/http"
)

func main() {
	println("Starting...")

	var errors []error
	svc := goVirtualHost.NewService()

	// virtual host: localhost
	errors = svc.Add(&goVirtualHost.HostInfo{
		Listens:   []string{":8080"},
		HostNames: []string{"localhost"},
		Handler:   http.FileServer(http.Dir(".")),
	})
	if len(errors)>0 {
		fmt.Println(errors)
	}

	// virtual host: default
	errors = svc.Add(&goVirtualHost.HostInfo{
		Listens: []string{":8080"},
		Handler: http.FileServer(http.Dir("/tmp")),
	})
	if len(errors)>0 {
		fmt.Println(errors)
	}

	// start server
	svc.Open()
}
