package handler

import (
	"context"

	"transfer_system/biz/model"
	"transfer_system/logs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateAccount .
// @router /accounts [POST]
func (a *App) CreateAccount(ctx context.Context, c *app.RequestContext) {
	var req model.NewAccount
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := a.CreateAccountResp(ctx, &req)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

func (a *App) CreateAccountResp(ctx context.Context, params *model.NewAccount) (*model.EmptyResponse, error) {
	_, err := a.CreateAccountService.Create(ctx, params)
	if err != nil {
		logs.CtxError(ctx, "failed to create account", err, logs.Fields{
			"params": params,
		})
		return nil, err
	}

	logs.CtxInfo(ctx, "successfully created account", logs.Fields{
		"params": params,
	})
	return &model.EmptyResponse{}, nil
}
