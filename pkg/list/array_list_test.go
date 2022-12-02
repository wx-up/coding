package list

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	slice := []int64{1, 2, 3}
	fmt.Println(slice[len(slice):])
}
