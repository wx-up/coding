package esacpe

/*
	go build -gcflags="-m" channel.go
	channel 传递指针会发生逃逸
*/

type Item struct{}

func Channel() {
	ch := make(chan *Item, 2)
	item := &Item{}
	ch <- item
}
