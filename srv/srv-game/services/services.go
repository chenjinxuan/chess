package services

import (
	"strings"
	"sync"
	"sync/atomic"

	. "chess/common/consul"
	"chess/common/log"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// a single connection
type client struct {
	id      string
	address string
	port    int
	conn    *grpc.ClientConn
}

// a kind of service
type service struct {
	clients []client
	idx     uint32 // for round-robin purpose
}

// all services
type service_pool struct {
	services       map[string]*service
	names          map[string]bool
	names_provided bool
	callbacks      map[string][]chan string
	mu             sync.RWMutex
}

var (
	_default_pool service_pool
	once          sync.Once
)

// 传入用到的服务
func Init(services []string) {
	once.Do(func() { _default_pool.init(services) })
}

func (p *service_pool) init(services []string) {

	// init
	p.services = make(map[string]*service)
	p.names = make(map[string]bool)

	// names init
	names := services // c.StringSlice("services")
	if len(names) > 0 {
		p.names_provided = true
	}

	log.Info("all service names:", names)
	for _, v := range names {
		p.names[strings.TrimSpace(v)] = true
	}

	// start connection
	p.connect_all()
}

// connect to all services
func (p *service_pool) connect_all() {
	// get services
	consulServices, err := ConsulClient.Services()
	if err != nil {
		log.Error(err)
		return
	}

	for _, consulService := range consulServices {
		p.add_service(consulService.ID, consulService.Service, consulService.Address, consulService.Port)
	}
	log.Info("services add complete")

	p.watcher()
}

// watcher for data change in etcd directory
func (p *service_pool) watcher() {
	for service_name := range p.services {
		go func(service string) {
			ConsulClient.ServiceWatch(service, func(idx uint64, raw interface{}) {
				if raw == nil {
					return
				}
				v, ok := raw.([]*consulapi.ServiceEntry)
				if !ok || len(v) == 0 {
					fmt.Println("consul return invalid")
				} else {
					for _, c := range p.services[service].clients {
						exists := false
						for _, s := range v {
							if s.Service.ID == c.id {
								if s.Service.Address != c.address || s.Service.Port != c.port { // 更新client
									p.update_service(s.Service.ID, service, s.Service.Address, s.Service.Port)
								}
								exists = true
								break
							}
						}
						if !exists { // 删除client
							p.remove_service(c.id, service)
						}
					}
					for _, s := range v {
						log.Infof("Watching service(%s) get:%+v", service, *s.Service)
						exists := false
						for _, c := range p.services[service].clients {
							if s.Service.ID == c.id {
								exists = true
								break
							}
						}
						if !exists { // 新增client
							p.add_service(s.Service.ID, service, s.Service.Address, s.Service.Port)
						}
					}
				}
			})
		}(service_name)
	}
}

// add a service
func (p *service_pool) add_service(id, name, address, port string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// name check
	service_name := name
	if p.names_provided && !p.names[service_name] {
		return
	}

	// try new service kind init
	if p.services[service_name] == nil {
		p.services[service_name] = &service{}
	}

	// create service connection
	service := p.services[service_name]
	target := fmt.Sprintf("%s:%d", address, port)
	if conn, err := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		service.clients = append(service.clients, client{id, address, port, conn})
		log.Info("service added:", id, "-->", target)
		//for k := range p.callbacks[service_name] {
		//	select {
		//	case p.callbacks[service_name][k] <- id:
		//	default:
		//	}
		//}
	} else {
		log.Info("did not connect:", id, "-->", target, "error:", err)
	}
}

// add a service
func (p *service_pool) update_service(id, name, address, port string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// name check
	service_name := name
	if p.names_provided && !p.names[service_name] {
		return
	}

	// try new service kind init
	if p.services[service_name] == nil {
		p.add_service(id, name, address, port)
		return
	}

	service := p.services[service_name]
	target := fmt.Sprintf("%s:%d", address, port)
	for k := range service.clients {
		if service.clients[k].id == id { // update
			// close old conn
			service.clients[k].conn.Close()

			oldtarget := fmt.Sprintf("%s:%d", service.clients[k].address, service.clients[k].port)
			if conn, err := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
				service.clients[k].address = address
				service.clients[k].port = port
				service.clients[k].conn = conn
				log.Infof("service(%s) update: %s --> %s", oldtarget, target)
				//for k := range p.callbacks[service_name] {
				//	select {
				//	case p.callbacks[service_name][k] <- id:
				//	default:
				//	}
				//}
			} else {
				log.Info("update service fail: did not connect:", id, "-->", target, "error:", err)
			}
			return
		}
	}
}

// remove a service
func (p *service_pool) remove_service(id, name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// name check
	service_name := name
	if p.names_provided && !p.names[service_name] {
		return
	}

	// check service kind
	service := p.services[service_name]
	if service == nil {
		log.Warnf("no such service:", service_name)
		return
	}

	// remove a service
	for k := range service.clients {
		if service.clients[k].id == id { // deletion
			service.clients[k].conn.Close()
			service.clients = append(service.clients[:k], service.clients[k+1:]...)
			log.Debug("service removed:", id)
			return
		}
	}
}

// provide a specific key for a service, eg: room-1
func (p *service_pool) get_service_with_id(id, name string) *grpc.ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	service := p.services[name]
	if service == nil {
		return nil
	}
	if len(service.clients) == 0 {
		return nil
	}

	// loop find a service with id
	for k := range service.clients {
		if service.clients[k].id == id {
			return service.clients[k].conn
		}
	}

	return nil
}

// get a service in round-robin style
// especially useful for load-balance with state-less services
func (p *service_pool) get_service(name string) (conn *grpc.ClientConn, id string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	service := p.services[name]
	if service == nil {
		return nil, ""
	}

	if len(service.clients) == 0 {
		return nil, ""
	}

	// get a service in round-robind style,
	idx := int(atomic.AddUint32(&service.idx, 1)) % len(service.clients)
	return service.clients[idx].conn, service.clients[idx].id
}

func (p *service_pool) register_callback(name string, callback chan string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.callbacks == nil {
		p.callbacks = make(map[string][]chan string)
	}

	p.callbacks[name] = append(p.callbacks[name], callback)
	if s, ok := p.services[name]; ok {
		for k := range s.clients {
			callback <- s.clients[k].id
		}
	}
	log.Info("register callback on:", name)
}

func GetService(name string) *grpc.ClientConn {
	conn, _ := _default_pool.get_service(name)
	return conn
}

func GetService2(name string) (*grpc.ClientConn, string) {
	conn, key := _default_pool.get_service(name)
	return conn, key
}

func GetServiceWithId(id, name string) *grpc.ClientConn {
	return _default_pool.get_service_with_id(id, name)
}

func RegisterCallback(name string, callback chan string) {
	_default_pool.register_callback(name, callback)
}
