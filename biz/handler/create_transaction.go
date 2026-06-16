package handler

import (
	"context"
	"transfer_system/biz/model"
	"transfer_system/logs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateTransaction .
// @router /transactions [POST]
func (a *App) CreateTransaction(ctx context.Context, c *app.RequestContext) {
	var req model.Transaction
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := a.CreateTransactionResp(ctx, &req)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

func (a *App) CreateTransactionResp(ctx context.Context, params *model.Transaction) (*model.EmptyResponse, error) {
	_, err := a.CreateTransactionService.Create(ctx, params)
	if err != nil {
		logs.CtxError(ctx, "failed to create transaction", err, logs.Fields{
			"params": params,
		})
		return nil, err
	}

	logs.CtxInfo(ctx, "successfully created transaction", logs.Fields{
		"params": params,
	})
	return &model.EmptyResponse{}, nil
}
