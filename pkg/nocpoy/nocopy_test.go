package nocpoy

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	u1 := Url{}
	u2 := u1
	fmt.Println(u2)
}
