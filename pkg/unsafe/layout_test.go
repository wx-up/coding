package unsafe

import "testing"

func TestPrintFieldOffset(t *testing.T) {
	PrintFieldOffset(struct {
		Name string
		Age  int64
	}{
		Name: "Bob",
		Age:  12,
	})
}
