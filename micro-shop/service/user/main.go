package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/google/uuid"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/wx-up/coding/micro-shop/service/user/global"
	"github.com/wx-up/coding/micro-shop/service/user/handler"
	"github.com/wx-up/coding/micro-shop/service/user/model"
	"github.com/wx-up/coding/micro-shop/service/user/proto"
	"google.golang.org/grpc"
)

func bootstrap() {
	// 初始化数据库
	global.InitDB()

	// migrate
	if err := global.DB().AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}

	// 初始化 consul
	global.InitConsul("127.0.0.1", 8500)
}

var (
	serverHost string
	serverPort int
)

func init() {
	flag.StringVar(&serverHost, "host", "0.0.0.0", "服务地址")
	flag.IntVar(&serverPort, "port", 8080, "服务端口")
	flag.Parse()
}

func main() {
	// 初始化逻辑
	bootstrap()

	// 启动服务
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserService{})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		panic(err)
	}

	// 注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 注册服务
	err = global.Consul().Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: "user",
		// name 和 id 都相同，多次启动会覆盖
		// name 相同 id 不同，相当于启动了多个相同的 server 可用于负载均衡
		ID:      uuid.New().String(),
		Tags:    []string{"user", "shop"},
		Port:    serverPort,
		Address: serverHost, // 这里是服务发现时给到的 IP 地址，不能是 0.0.0.0/127.0.0.1 得是本机的一个 IP 地址
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serverHost, serverPort),
			Interval:                       "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		fmt.Printf("服务注册失败：%v\n", err)
	}

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
