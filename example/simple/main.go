package main

import (
	"fmt"
	"mjpclab.dev/goVirtualHost"
	"net/http"
)

func main() {
	println("Starting...")

	var errors []error
	svc := goVirtualHost.NewService()

	// virtual host: localhost
	errors, _ = svc.Add(&goVirtualHost.HostInfo{
		Listens:   []string{":8080"},
		HostNames: []string{"localhost"},
		Handler:   http.FileServer(http.Dir(".")),
	})
	if len(errors) > 0 {
		fmt.Println(errors)
	}

	// virtual host: default
	errors, _ = svc.Add(&goVirtualHost.HostInfo{
		Listens: []string{":8080"},
		Handler: http.FileServer(http.Dir("/tmp")),
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
