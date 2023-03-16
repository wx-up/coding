//go:build e2e

package micro

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(":8081")
	require.NoError(t, err)
	srv := &UserService{}
	err = InitClientProxy(srv, client)
	require.NoError(t, err)
	resp, err := srv.GetById(context.Background(), &GetByIdReq{Id: 10})
	require.NoError(t, err)
	log.Println(resp)
}
