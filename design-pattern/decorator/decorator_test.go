package decorator

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	sc := &SafeCache{
		Cache: &memoryCache{
			values: make(map[string]string),
		},
	}
	fmt.Println(sc)
}
