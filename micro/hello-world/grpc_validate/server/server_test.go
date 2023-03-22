package main

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
)

func Test(t *testing.T) {
	fmt.Println(codes.InvalidArgument)
}
