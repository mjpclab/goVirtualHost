package main

import (
	goVirtualHost "../.."
	"fmt"
	"net/http"
)

func main() {
	println("Starting...")

	var err error
	svc := goVirtualHost.NewService()

	// virtual host: localhost
	err = svc.Add(&goVirtualHost.HostInfo{
		Listen:    ":8080",
		HostNames: []string{"localhost"},
		Handler:   http.FileServer(http.Dir(".")),
	})
	if err != nil {
		fmt.Println(err)
	}

	// virtual host: default
	err = svc.Add(&goVirtualHost.HostInfo{
		Listen:  ":8080",
		Handler: http.FileServer(http.Dir("/tmp")),
	})
	if err != nil {
		fmt.Println(err)
	}

	// start server
	svc.Open()
}
