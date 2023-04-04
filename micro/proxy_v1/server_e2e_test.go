//go:build e2e

package proxy_v1

import "testing"

func TestServer_Start(t *testing.T) {
	server := &Server{}
	server.Start(":8081")
}
