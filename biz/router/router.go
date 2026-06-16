package router

import (
	"context"

	"transfer_system/biz/handler"

	hertzapp "github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

/*
Register registers routes for the transfer system API.
*/
func Register(r *server.Hertz, handlerApp *handler.App) {
	root := r.Group("/")
	{
		root.GET("/ping", func(ctx context.Context, c *hertzapp.RequestContext) {
			c.String(consts.StatusOK, "pong")
		})
		root.GET("/accounts/:account_id", handlerApp.GetAccount)
		root.POST("/accounts", handlerApp.CreateAccount)
		root.POST("/transactions", handlerApp.CreateTransaction)
	}
}
