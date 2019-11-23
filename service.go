package goVirtualHost

import "sync"

func NewService() *Service {
	service := &Service{
		listeners: []*listener{},
		servers:   []*server{},
		vhosts:    []*vhost{},
	}

	return service
}

func (svc *Service) addParam(param *param) {
	// params
	svc.params = append(svc.params, param)

	// listeners, servers
	var listener *listener
	var server *server

	listener = svc.listeners.find(param.proto, param.addr)
	if listener != nil {
		server = listener.server
	} else {
		server = newServer(param.useTLS)
		listener = newListener(param.proto, param.addr)
		listener.server = server

		svc.listeners = append(svc.listeners, listener)
		svc.servers = append(svc.servers, server)
	}

	// vhost
	vhost := newVhost(param.cert, param.hostNames, param.handler)
	svc.vhosts = append(svc.vhosts, vhost)

	// server -> vhost
	server.vhosts = append(server.vhosts, vhost)
}

func (svc *Service) Add(host *HostInfo) []error {
	errors := []error{}

	newParams := host.toParams()
	for _, newParam := range newParams {
		err := svc.params.validate(newParam)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		svc.addParam(newParam)
	}

	return errors
}

func (svc *Service) Open() (errors []error) {
	for _, s := range svc.servers {
		s.updateDefaultVhost()
		s.updateHttpServerTLSConfig()
		s.updateHttpServerHandler()
	}

	errors = []error{}
	chErr := make(chan error)

	wgOpen := sync.WaitGroup{}
	for _, listener := range svc.listeners {
		wgOpen.Add(1)
		l := listener
		go func() {
			defer wgOpen.Done()

			// start net listener
			err := l.open()
			if err != nil {
				chErr <- err
				return
			}

			// start http serve
			err = l.server.open(l)
			if err != nil {
				chErr <- err
				return
			}

		}()
	}

	wgErr := sync.WaitGroup{}
	go func() {
		wgErr.Add(1)
		for err := range chErr {
			errors = append(errors, err)
		}
		wgErr.Done()
	}()

	wgOpen.Wait()
	close(chErr)
	wgErr.Wait()

	return
}

func (svc *Service) Close() {
	wg := sync.WaitGroup{}
	for _, listener := range svc.listeners {
		wg.Add(1)
		l := listener
		go func() {
			l.server.close()
			l.close()
			wg.Done()
		}()
	}
	wg.Wait()
}
