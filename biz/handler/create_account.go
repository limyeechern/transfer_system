package handler

import (
	"context"

	"transfer_system/biz/model"
	"transfer_system/biz/service/create_account"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/sirupsen/logrus"
)

var (
	CreateAccountService create_account.CreateAccountService
)

func init() {
	CreateAccountService = create_account.NewCreateAccountService()
}

// CreateAccount .
// @router /accounts [POST]
func CreateAccount(ctx context.Context, c *app.RequestContext) {
	var req model.NewAccount
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := CreateAccountResp(ctx, &req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

func CreateAccountResp(ctx context.Context, params *model.NewAccount) (*model.EmptyResponse, error) {
	_, err := CreateAccountService.Create(ctx, params)
	if err != nil {
		logrus.WithContext(ctx).WithField("params", params).Error("failed to create account")
		return nil, err
	}

	logrus.WithContext(ctx).Info("successfully created account")
	return &model.EmptyResponse{}, nil
}
