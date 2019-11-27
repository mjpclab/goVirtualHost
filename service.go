package goVirtualHost

import "sync"

func NewService() *Service {
	service := &Service{
		listeners: listeners{},
		servers:   servers{},
		vhosts:    vhosts{},
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

func (svc *Service) Add(info *HostInfo) (errs []error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	newParams := info.toParams()
	for _, newParam := range newParams {
		err := svc.params.validate(newParam)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		svc.addParam(newParam)
	}

	return errs
}

func (svc *Service) openListeners() (errs []error) {
	chListenErr := make(chan error)

	go func() {
		wg := sync.WaitGroup{}
		for _, listener := range svc.listeners {
			wg.Add(1)
			l := listener
			go func() {
				err := l.open()
				if err != nil {
					chListenErr <- err
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(chListenErr)
	}()

	for err := range chListenErr {
		errs = append(errs, err)
	}

	return
}

func (svc *Service) openServers(cbAllArranged func()) (errs []error) {
	chServeErr := make(chan error)

	go func() {
		wg := sync.WaitGroup{}
		for _, listener := range svc.listeners {
			wg.Add(1)
			l := listener
			go func() {
				err := l.server.open(l)
				if err != nil {
					chServeErr <- err
				}
				wg.Done()
			}()
		}
		cbAllArranged()
		wg.Wait()
		close(chServeErr)
	}()

	for err := range chServeErr {
		errs = append(errs, err)
	}

	return
}

func (svc *Service) Open() (errs []error) {
	svc.mu.Lock()

	for _, s := range svc.servers {
		s.updateDefaultVhost()
		s.updateHttpServerTLSConfig()
		s.updateHttpServerHandler()
	}

	errs = svc.openListeners()
	if len(errs) > 0 {
		svc.mu.Unlock()
		return errs
	}

	errs = svc.openServers(svc.mu.Unlock)
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (svc *Service) Close() {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	wg := sync.WaitGroup{}
	for _, listener := range svc.listeners {
		wg.Add(1)
		l := listener
		go func() {
			if l.server != nil {
				l.server.close()
			}
			l.close()
			wg.Done()
		}()
	}
	wg.Wait()
}
