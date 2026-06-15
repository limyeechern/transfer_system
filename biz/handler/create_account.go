package handler

import (
	"context"

	"transfer_system/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/sirupsen/logrus"
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
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

func (a *App) CreateAccountResp(ctx context.Context, params *model.NewAccount) (*model.EmptyResponse, error) {
	_, err := a.CreateAccountService.Create(ctx, params)
	if err != nil {
		logrus.WithContext(ctx).WithField("params", params).Error("failed to create account")
		return nil, err
	}

	logrus.WithContext(ctx).Info("successfully created account")
	return &model.EmptyResponse{}, nil
}
