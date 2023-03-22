package main

import (
	"fmt"

	"github.com/wx-up/coding/micro/hello-world/grpc_validate/proto"
)

func main() {
	obj := new(proto.User)
	fmt.Println(obj.Validate())
}
