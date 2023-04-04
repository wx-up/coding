package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wx-up/coding/micro-shop/web/pkg/consul/load_balancing"

	"github.com/wx-up/coding/micro-shop/web/pkg"

	"github.com/wx-up/coding/micro-shop/web/pkg/client"

	"github.com/wx-up/coding/micro-shop/web/config"

	"github.com/gin-gonic/gin"
	"github.com/wx-up/coding/micro-shop/web/router"
)

var (
	host string
	port int
)

func init() {
	time.AfterFunc(time.Second, func() {
	})
	flag.IntVar(&port, "port", 8088, "web 服务端口")
	flag.StringVar(&host, "host", "0.0.0.0", "web 服务地址")
	flag.Parse()
}

func InitClient() {
	conn, err := load_balancing.GetConnByServiceName(config.Config().UserServer.Name, func(conf *load_balancing.Config) {
		consulConf := config.Config().Consul
		conf.SetPort(consulConf.Port).SetHost(consulConf.Host)
	})
	if err != nil {
		panic(err)
	}
	client.Register(config.Config().UserServer.Name, conn)
}

func main() {
	// 初始化日志
	config.InitConfig("micro-shop/web/config.yml")

	// 初始化路由
	engine := gin.New()
	// gin.SetMode(gin.ReleaseMode)
	router.RegisterRouter(engine)

	// 初始化服务连接
	InitClient()

	// 初始化 consul
	//consulClient := consul.New(config.Config().Consul.Host, config.Config().Consul.Port)
	//
	//// 获取用户服务的地址和端口
	//infos, err := consulClient.GetServiceInfos(config.Config().UserServer.Name)
	//if err != nil || len(infos) <= 0 {
	//	fmt.Printf("获取服务信息失败：%v", err)
	//	return
	//}
	//info := infos[0]
	//client.InitClient(info.Host, info.Port)

	if config.Config().Env == "prod" {
		freePort, _ := pkg.GetFreePort()
		if freePort > 0 {
			port = freePort
		}
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: engine,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("服务器异常退出：%v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 关闭服务
	client.Close()

	fmt.Println("服务器成功退出")
}
