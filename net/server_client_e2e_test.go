package net

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	go func() {
		srv := NewServer(":8080")
		err := srv.Listen()
		require.Nil(t, err)
	}()

	time.Sleep(time.Second * 3)

	client := NewClient(":8080")
	resp, err := client.Send("星期三")
	require.Nil(t, err)
	assert.Equal(t, "星期三, from response", resp)
}
