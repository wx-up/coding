package builder

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	u := NewUserBuilder().BuildName("wx").Build()
	fmt.Printf("%+v", u)
}
