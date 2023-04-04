package v2

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	ch := make(chan int, 10)
	fmt.Println(cap(ch))
}
