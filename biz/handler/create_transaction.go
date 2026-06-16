package handler

import (
	"context"
	"transfer_system/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/sirupsen/logrus"
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
		logrus.WithContext(ctx).WithField("params", params).Error("failed to create transaction")
		return nil, err
	}

	logrus.WithContext(ctx).Info("successfully created transaction")
	return &model.EmptyResponse{}, nil
}
