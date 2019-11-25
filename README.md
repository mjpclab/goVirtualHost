# goVirtualHost
goVirtualHost: An easy way to setup HTTP virtual hosts.

Minimal required Go version is 1.9.

# Quick Example
Two virtual hosts both listen on :8080, but with different hostname, serve for different directories:
```go
svc := goVirtualHost.NewService()

// virtual host: localhost
svc.Add(&goVirtualHost.HostInfo{
    Listens:   []string{":8080"},
    HostNames: []string{"localhost"},
    Handler:   http.FileServer(http.Dir(".")),
})

// virtual host: default host
svc.Add(&goVirtualHost.HostInfo{
    Listens: []string{":8080"},
    Handler: http.FileServer(http.Dir("/tmp")),
})

// start server
svc.Open()
```

# NewService(&HostInfo) *Service
`NewService` returns a service that manages multiple virtual hosts.

# HostInfo
the `HostInfo` is the initial virtual host information, with the properties:

## Listens []string
IP and/or port the server listens on, e.g. ":80" or "127.0.0.1:80".
if `Cert` is present, Serve for TLS HTTP, otherwise Serve for plain HTTP.
If port is not specified, use "80" as default for Plain HTTP mode, "443" for TLS mode.
If value contains "/" then treat it as a unix socket file.

## ListensPlain []string
IP and/or port the server listens on, e.g. ":80" or "127.0.0.1:80".
Serve for plain HTTP.
If port is not specified, use "80" as default.
If value contains "/" then treat it as a unix socket file.

## ListensTLS []string
IP and/or port the server listens on, e.g. ":443" or "127.0.0.1:443".
Serve for TLS HTTP.
If port is not specified, use "443" as default.
If value contains "/" then treat it as a unix socket file.

## Cert *tls.Certificate
TLS Certificate supplied for TLS mode. To load from external PEM files, use:
```go
func getCert(certFile, keyFile string) *tls.Certificate {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err !=nil {
        return nil
    }
    return &cert
}
```

## HostNames []string
Specify hostnames associated with the virtual host.
If hostname starts with ".", treat it as a suffix, to match all levels of sub domains, e.g. ".example.com".
If request host name does not match any virtual host,
server will try to use first virtual host that has no hostname,
otherwise use the first virtual host.

## Handler http.Handler
`http.Handler` to handle requests. Could be a `http.ServeMux`.

# *Service.Open() []error
Start listening on network ports, and serve for http requests. The method will not return until all servers are closed.
e.g. call `Close` method on another goroutine.

# *Service.Close()
Stop serving. To restart serving, a new `Service` must be created.

# Architecture & Internals
```
Service
    |  manages
    v
    +---------+---------+---------+---------+
    | handler | handler | handler | handler |
    +---------+---------+---------+---------+
    |  vhost  |  vhost  |  vhost  |  vhost  |
    +---------+--+------+-----+---+---------+
    |   server   |   server   |   server    |
    +------------+------------+-------------+
    |  listener  |  listener  |  listener   |
    +------------+------------+-------------+
```

## listener
`listener` is a wrapper for `net.Listener`, which open ports or sockes and listen.

## server
`server` is a wrapper for `http.Server`. It's `handler` does not serve for end user,
but dispatching requests to related virtual host according to the Host header.

## vhost
`vhost` manages related hostnames, and hold `handler` to deal with requests dispatched from `server`.
