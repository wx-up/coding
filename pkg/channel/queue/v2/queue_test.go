package v2

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	broker := &Broker{}
	broker.Subscribe(NewConsumerGroup("测试", 2))
	broker.Publish("哈哈")
	broker.Publish("你好")
	broker.Publish("你好23")
	broker.Publish("你好45")
	time.Sleep(time.Second * 10)
	broker.Publish("嘻嘻")
	broker.Publish("哈哈")
	time.Sleep(time.Second * 2)
}
