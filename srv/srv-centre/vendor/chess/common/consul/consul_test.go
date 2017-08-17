package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"
)

func TestConsulCliWrap_RegisterService(t *testing.T) {
	err := InitConsulClientViaEnv()
	if err != nil {
		t.Fatal(err)
	}

	err = ConsulClient.RegisterService(&api.AgentServiceRegistration{
		ID:                "room-1",
		Name:              "room",
		Tags:              []string{"master"},
		Port:              8000,
		Address:           "127.0.0.1",
		EnableTagOverride: true,
	})

	err = ConsulClient.RegisterService(&api.AgentServiceRegistration{
		ID:                "room-2",
		Name:              "room",
		Tags:              []string{"master"},
		Port:              8001,
		Address:           "127.0.0.1",
		EnableTagOverride: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestConsulCliWrap_DeregisterService(t *testing.T) {
	err := ConsulClient.DeregisterService("room-2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestConsulCliWrap_Services(t *testing.T) {
	services, err := ConsulClient.Services()
	if err != nil {
		t.Fatal(err)
	}

	for id, service := range services {
		t.Logf("ID:%s  Service:%+v", id, service)
	}
}

func TestConsulCliWrap_ServiceWatch(t *testing.T) {
	go func() {
		time.Sleep(5 * time.Second)
		ConsulClient.RegisterService(&api.AgentServiceRegistration{
			ID:                "room-3",
			Name:              "room",
			Tags:              []string{"master"},
			Port:              8000,
			Address:           "127.0.0.1",
			EnableTagOverride: true,
		})
		time.Sleep(5 * time.Second)
		ConsulClient.DeregisterService("room-3")
	}()

	err := ConsulClient.ServiceWatch("room", func(idx uint64, raw interface{}) {
		//fmt.Println(raw)
		if raw == nil {
			return
		}
		v, ok := raw.([]*api.ServiceEntry)
		if !ok || len(v) == 0 {
			fmt.Println("consul return invalid")
		} else {
			for _, s := range v {
				fmt.Println(*s.Service)
			}
		}
		fmt.Println("--------")
	})
	if err != nil {
		t.Fatal(err)
	}
}
