package main

import (
	"github.com/hashicorp/consul/api"
)

func Register(name string, address string, port int, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id + "2",
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: address,

		Check: &api.AgentServiceCheck{
			Interval:                       "10s",
			Timeout:                        "5s",
			HTTP:                           "http://localhost:8081/health-check",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func GetServicesByFilter(filter string) (map[string]*api.AgentService, error) {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return client.Agent().ServicesWithFilter(filter)
}

func main() {
	err := Register("web-server", "127.0.0.1", 8081, []string{"gin", "web"}, "web-server")
	if err != nil {
		panic(err)
	}
}
