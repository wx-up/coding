package main

import (
	"fmt"
	"testing"

	_ "github.com/mbobakov/grpc-consul-resolver"

	"google.golang.org/grpc"
)

func Test(t *testing.T) {
	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/whoami?wait=14s&tag=manual",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
}

func TestChannel(t *testing.T) {
	b := struct{}{}

	c := struct{}{}

	fmt.Println(b == c)
}
