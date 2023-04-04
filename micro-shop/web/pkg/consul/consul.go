package consul

import (
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	port int
	host string

	instance *api.Client
}

func New(host string, port int) *Client {
	client := &Client{
		host: host,
		port: port,
	}
	client.instance = initConsul(host, port)
	return client
}

type ServiceInfo struct {
	Port int
	Host string
}

func (c *Client) GetServiceInfos(name string) ([]ServiceInfo, error) {
	res, err := c.instance.Agent().ServicesWithFilter(fmt.Sprintf("Service == %s", name))
	if err != nil {
		return nil, err
	}
	if len(res) <= 0 {
		return nil, errors.New("服务不存在")
	}
	infos := make([]ServiceInfo, 0, len(res))

	for _, v := range res {
		infos = append(infos, ServiceInfo{
			Port: v.Port,
			Host: v.Address,
		})
	}
	return infos, nil
}

func (c *Client) Close() {
	if c.instance != nil {
	}
}

func initConsul(host string, port int) *api.Client {
	conf := api.DefaultConfig()
	conf.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	return client
}
