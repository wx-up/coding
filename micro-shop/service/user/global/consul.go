package global

import (
	"fmt"
	"sync"

	"github.com/hashicorp/consul/api"
)

var (
	instance *api.Client
	onceInit sync.Once
)

func Consul() *api.Client {
	return instance
}

func InitConsul(host string, port int) {
	onceInit.Do(func() {
		initConsul(host, port)
	})
}

func initConsul(host string, port int) {
	conf := api.DefaultConfig()
	conf.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	instance = client
}
