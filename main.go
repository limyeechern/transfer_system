package main

import (
	"fmt"

	"transfer_system/biz/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	addr := ":8080"
	h := server.Default(server.WithHostPorts(addr))
	router.Register(h)

	fmt.Printf("server listening on %s\n", addr)
	h.Spin()
}
