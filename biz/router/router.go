package router

import (
	"context"

	"transfer_system/biz/handler"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

/*
Register registers routes for the transfer system API.
*/
func Register(r *server.Hertz) {
	root := r.Group("/")
	{
		root.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
			c.String(consts.StatusOK, "pong")
		})
		root.POST("/accounts", handler.CreateAccount)
	}
}
