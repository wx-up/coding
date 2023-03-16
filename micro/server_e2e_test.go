//go:build e2e

package micro

import "testing"

func TestServer_Start(t *testing.T) {
	server := &Server{}
	server.Start(":8081")
}
