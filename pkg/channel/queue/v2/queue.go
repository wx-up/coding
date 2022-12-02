package v2

import (
	"context"
	"fmt"
	"time"
)

type Broker struct {
	groups []*ConsumerGroup
}

func (b *Broker) Subscribe(c *ConsumerGroup) {
	b.groups = append(b.groups, c)
}

func (b *Broker) Publish(message string) {
	for _, g := range b.groups {
		g.consume(message)
	}
}

// ConsumerGroup 消费组
type ConsumerGroup struct {
	ch    chan string
	name  string
	maxGo chan struct{}
}

func NewConsumerGroup(name string, maxCount int64) *ConsumerGroup {
	c := &ConsumerGroup{
		name:  name,
		ch:    make(chan string, maxCount),
		maxGo: make(chan struct{}, maxCount), // 控制最大 goroutine 数量
	}
	return c
}

func (c *ConsumerGroup) consume(message string) {
	select {
	case c.ch <- message:
		select {
		case c.maxGo <- struct{}{}:
			fmt.Println("启动一个 goroutine")
			go c.singleHande()
		default: // 不阻塞
		}
	}
}

func (c *ConsumerGroup) singleHande() {
	defer func() {
		<-c.maxGo
	}()
	for {
		// 控制长时间获取不到任务之后，就退出
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*4)
		select {
		case task := <-c.ch:
			fmt.Println(task)
			cancelFunc()
		case <-ctx.Done():
			cancelFunc()
			fmt.Println(" goroutine 长时间没有消费，退出了")
			return
		}
	}
}
