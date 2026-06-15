package main

import (
	"context"
	"fmt"
	"log"

	bootstrap "transfer_system/biz/app"
	"transfer_system/biz/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	ctx := context.Background()
	app, cleanup, err := bootstrap.New(ctx, bootstrap.ConfigFromEnv())
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	defer cleanup()

	addr := ":8080"
	h := server.Default(server.WithHostPorts(addr))
	router.Register(h, app)

	fmt.Printf("server listening on %s\n", addr)
	h.Spin()
}
