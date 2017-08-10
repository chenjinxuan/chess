package services

import (
	"strings"
	"sync"
	//"sync/atomic"
	. "chess/common/consul"
	"chess/common/log"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"time"
	"sync/atomic"
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
	clients map[string]*client   // service-id -> client
	ids []string // for round-robin purpose
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

// 发现服务 - 传入用到的服务名
func Discover(services []string) {
	once.Do(func() { _default_pool.init(services) })
}

// 注册服务
func Register(id, name, address string, port, checkPort int, tags []string) error {
	return ConsulClient.RegisterService(&consulapi.AgentServiceRegistration{
		ID:                id,
		Name:              name,
		Tags:              tags,
		Port:              port,
		Address:           address,
		EnableTagOverride: true,
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d%s", address, checkPort, "/check"),
			Timeout:                        "3s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
		},
	})
}

// 注销服务
func Deregister(id string) error {
	return ConsulClient.DeregisterService(id)
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
		name := strings.TrimSpace(v)
		p.names[name] = true

		consulServices, err := ConsulClient.Service(name)
		if err != nil {
			log.Errorf("ConsulClient.Service(%s) Error: %s", name, err)
			continue
		}

		for _, consulService := range consulServices {
			p.add_service(consulService.ServiceID, consulService.ServiceName, consulService.ServiceAddress, consulService.ServicePort)
		}

	}
	log.Info("services add complete")
	// start connection
	p.watcher()
}

// watcher for data change in etcd directory
func (p *service_pool) watcher() {
	for service_name := range p.names {
		go func(service string) {
			ConsulClient.ServiceWatch(service, func(idx uint64, raw interface{}) {
				if raw == nil {
					return
				}
				v, ok := raw.([]*consulapi.ServiceEntry)
				if !ok {
					fmt.Println("consul return invalid")
				} else {
					if p.services[service] != nil {
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
					} else {
						for _, s := range v {
							log.Infof("Watching service(%s) get:%+v", service, *s.Service)
							// 新增client
							p.add_service(s.Service.ID, service, s.Service.Address, s.Service.Port)
						}
					}

				}
			})
		}(service_name)
	}
}

// add a service
func (p *service_pool) add_service(id, name, address string, port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// name check
	service_name := name
	if p.names_provided && !p.names[service_name] {
		return
	}

	// try new service kind init
	if p.services[service_name] == nil {
		p.services[service_name] = &service{
			clients: make(map[string]*client),
		}
	}

	// create service connection
	service := p.services[service_name]
	if service.clients[id] != nil {
		return
	}

	target := fmt.Sprintf("%s:%d", address, port)
	if conn, err := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		service.clients[id] = &client{id, address, port, conn}
		service.ids = append(service.ids, id)
		log.Info("service added:", id, "-->", target)
		go func (n string, c *client){
			tg:= fmt.Sprintf("%s:%d", c.address,c.port)
			wantState := grpc.Shutdown
			if _, ok := assert_state(wantState, c.conn); ok { // 连接断开
				log.Info("service shutdown:", id, "-->", tg)
				p.remove_service(c.id, n)
			}
		}(service_name, service.clients[id])

	} else {
		log.Info("did not connect:", id, "-->", target, "error:", err)
	}
}

// add a service
func (p *service_pool) update_service(id, name, address string, port int) {
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
	if client, ok := service.clients[id]; ok {
		// close old conn
		client.conn.Close()

		oldtarget := fmt.Sprintf("%s:%d", client.address, client.port)
		if conn, err := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
			client.address = address
			client.port = port
			client.conn = conn
			log.Infof("service(%s) update: %s --> %s", oldtarget, target)
		} else {
			log.Info("update service fail: did not connect:", id, "-->", target, "error:", err)
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
	if client, ok := service.clients[id]; ok {
		log.Debug("service removed:", id, "-->", fmt.Sprintf("%s:%d", client.address, client.port))
		client.conn.Close()
		delete(service.clients, id)
		for k, v := range service.ids {
			if v == id {
				service.ids = append(service.ids[:k], service.ids[k+1:]...)
			}
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
	client, ok := service.clients[id]
	if ok && client.conn.GetState() != grpc.Shutdown {
		return client.conn
	}

	return nil
}

func (p *service_pool) get_service_lb(name string) (conn *grpc.ClientConn, id string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.get_service(name)
}

// get a service in round-robin style
// especially useful for load-balance with state-less services
func (p *service_pool) get_service(name string) (conn *grpc.ClientConn, id string) {
	// check existence
	service := p.services[name]
	if service == nil {
		return nil, ""
	}

	if len(service.clients) == 0 {
		return nil, ""
	}

	// get a service in round-robind style,
	idx := int(atomic.AddUint32(&service.idx, 1)) % len(service.ids)
	id = service.ids[idx]
	if service.clients[id].conn.GetState() == grpc.Shutdown {
		p.remove_service(id, name)
		return p.get_service(name)
	}

	return service.clients[id].conn, id
}

func GetService(name string) *grpc.ClientConn {
	conn, _ := _default_pool.get_service_lb(name)
	return conn
}

func GetService2(name string) (*grpc.ClientConn, string) {
	conn, key := _default_pool.get_service_lb(name)
	return conn, key
}

func GetServiceWithId(id, name string) *grpc.ClientConn {
	return _default_pool.get_service_with_id(id, name)
}

func assert_state(wantState grpc.ConnectivityState, cc *grpc.ClientConn) (grpc.ConnectivityState, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var state grpc.ConnectivityState
	for state = cc.GetState(); state != wantState && cc.WaitForStateChange(ctx, state); state = cc.GetState() {
	}
	return state, state == wantState
}