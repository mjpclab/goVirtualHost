# goVirtualHost
goVirtualHost: An easy way to setup HTTP virtual hosts.

Minimal required Go version is 1.14.

# Quick Example
Two virtual hosts listen on both :8080, but with different hostname, serve for different directories:
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

# NewService() *Service
`NewService` returns a service instance that manages multiple virtual hosts.

# (*Service) Add(*HostInfo) []error
Adding a new virtual host information to the `Service`.
You can use `errors.Is()` to test possible errors:

- `CertificateNotFound`

Intent to work on TLS mode, but certificate is not provided.

- `ConflictIPAddress`

One Virtual host tries to listen on a specific IP address and port,
while another virtual host tries to listen on a wildcard IP(e.g. "0.0.0.0") of same port.

- `ConflictTLSMode`

For a specific listening endpoint(IP:port or socket),
one virtual host works on plain mode,
while another virtual host works on TLS mode.

- `DuplicatedAddressHostname`

Two virtual hosts listen on same endpoint, they use the same hostname.

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
TLS Certificate supplied for TLS mode. A helper function can be used to load from external PEM files:
```go
func LoadCertificate(certFile, keyFile string) (*tls.Certificate, error)
```

## HostNames []string
Specify hostnames associated with the virtual host.
If hostname starts with ".", treat it as a suffix, to match all levels of sub domains, e.g. ".example.com".
If hostname ends with ".", treat it as a prefix, to match all levels of suffix domains, e.g. "192.168.1.".
If request host name does not match any virtual host,
server will try to use first virtual host that has no hostname,
otherwise use the first virtual host.

## Handler http.Handler
`http.Handler` to handle requests.
Could be an instance of `http.ServeMux`, `httputil.ReverseProxy`, or any other type that implements `http.Handler`.

# (*Service) Open() []error
Start listening on network ports, and serve for http requests. The method will not return until all servers are closed.
e.g. call `Close` method on another goroutine.

# (*Service) Close()
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
