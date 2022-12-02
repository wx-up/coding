package _map

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	sm := &SafeMap[string, string]{
		values: make(map[string]string, 10),
	}
	fmt.Println(sm.LoadOrStore("name", "bob"))
}
