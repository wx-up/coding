package load_balancing

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	host   string
	port   int
	policy string
}

func (c *Config) SetHost(host string) *Config {
	c.host = host
	return c
}

func (c *Config) SetPort(port int) *Config {
	c.port = port
	return c
}

func (c *Config) SetPolicy(policy string) *Config {
	c.policy = policy
	return c
}

func GetConnByServiceName(name string, apply func(config *Config)) (*grpc.ClientConn, error) {
	conf := &Config{
		host:   "127.0.0.1",
		port:   8500,
		policy: "round_robin",
	}
	apply(conf)
	return grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=10s", conf.host, conf.port, name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, conf.policy)),
	)
}
