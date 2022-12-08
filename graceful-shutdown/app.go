package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Option func(*App)

func WithShutdownCallbacks(cbs []ShutdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
	}
}

// ShutdownCallback 采用 context.Context 来控制超时，而不是用 time.After
// 因为希望用户知道，他的回调必须要在一定时间内处理完毕，而且他必须显式处理超时错误
type ShutdownCallback func(ctx context.Context)

// StoreCacheToDBCallback 将 cache 刷新到数据库中（ 例子 ）
func StoreCacheToDBCallback(ctx context.Context) {
	done := make(chan struct{}, 1)

	// 业务逻辑
	go func() {
		log.Printf("正在刷新")
		time.Sleep(time.Second)
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Printf("处理成功")
	case <-ctx.Done():
		log.Printf("刷新到DB超时")
	}

	return
}

type App struct {
	// 可以开启多个服务
	servers []*Server

	// 优雅退出整个超时时间
	shutdownTimeout time.Duration

	// 优雅退出等待已有请求处理的时间
	waitTimeout time.Duration

	// 回调处理超时时间
	cbTimeout time.Duration

	// 回调
	cbs []ShutdownCallback
}

func NewApp(servers []*Server, opts ...Option) *App {
	app := &App{
		servers:         servers,
		shutdownTimeout: 30 * time.Second,
		waitTimeout:     15 * time.Second,
		cbTimeout:       3 * time.Second,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *App) StartAndServe() {
	// 启动服务
	for _, s := range app.servers {
		srv := s
		go func() {
			if err := srv.Start(); err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器 %s 已经关闭", srv.name)
				} else {
					log.Printf("服务器 %s 异常退出", srv.name)
				}
			}
		}()
	}

	// 监听退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	<-ch
	go func() {
		select {
		case <-ch:
			log.Printf("强制退出")
			os.Exit(1)
		case <-time.After(app.shutdownTimeout):
			log.Printf("超时强制退出")
			os.Exit(1)
		}
	}()

	// 退出
	app.shutdown()
}

func (app *App) shutdown() {
	// 拒绝请求，等待现有的请求处理完毕，如果是 HTTP 服务，它的 Shutdown 方法就包含了这两步
	// 如果是 RPC 服务的话，就需要自己处理这块
	for _, s := range app.servers {
		// 拒绝请求
		s.rejectReq()
	}

	// 等待现有请求处理完成
	// TODO：这里可以改造为实时统计正在处理的请求数量，为0时再进行下一步，而不是粗略的根据时间判断
	time.Sleep(app.waitTimeout)

	// 关闭服务
	var wg sync.WaitGroup
	for _, s := range app.servers {
		wg.Add(1)
		srv := s
		go func() {
			defer wg.Done()
			if err := srv.Stop(); err != nil {
				log.Print("关闭服务失败", err)
			}
		}()
	}
	wg.Wait()

	// 执行回调
	for _, cb := range app.cbs {
		wg.Add(1)
		c := cb
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), app.cbTimeout)
			defer cancel()
			c(ctx)
		}()
	}
	wg.Wait()

	// 释放应用资源
	app.close()
}

func (app *App) close() {
	// 释放应用可能的一些资源
	time.Sleep(time.Second)
	log.Print("应用关闭")
}
