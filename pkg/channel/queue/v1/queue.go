package v1

import "fmt"

type Broker struct {
	consumers []*ConsumerGroup
}

func (b *Broker) Subscribe(c *ConsumerGroup) {
	b.consumers = append(b.consumers, c)
}

func (b *Broker) Publish(message string) {
	for _, consumer := range b.consumers {
		consumer.ch <- message
	}
}

type ConsumerGroup struct {
	ch            chan string
	name          string
	consumerCount int
}

func NewConsumerGroup(name string, consumerCount int) *ConsumerGroup {
	c := &ConsumerGroup{
		ch:            make(chan string, consumerCount),
		name:          name,
		consumerCount: consumerCount,
	}
	c.start()
	return c
}

func (c *ConsumerGroup) start() {
	for i := 0; i < c.consumerCount; i++ {
		go c.singleHandle()
	}
}

func (c *ConsumerGroup) singleHandle() {
	for {
		select {
		case task := <-c.ch:
			fmt.Println(task)
		}
	}
}
