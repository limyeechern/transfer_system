package handler

import (
	"context"
	"strconv"

	"transfer_system/biz/apperror"
	"transfer_system/biz/model"
	"transfer_system/logs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetAccount .
// @router /accounts/:account_id [GET]
func (a *App) GetAccount(ctx context.Context, c *app.RequestContext) {
	accountID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
	if err != nil {
		c.String(consts.StatusBadRequest, apperror.ErrInvalidAccount.Error())
		return
	}

	resp, err := a.GetAccountResp(ctx, &model.GetAccount{AccountID: accountID})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

func (a *App) GetAccountResp(ctx context.Context, params *model.GetAccount) (*model.Account, error) {
	respData, err := a.GetAccountService.Read(ctx, params)
	if err != nil {
		logs.CtxError(ctx, "failed to get account", err, logs.Fields{
			"params": params,
		})
		return nil, err
	}

	logs.CtxInfo(ctx, "successfully got account", logs.Fields{
		"account_id": respData.AccountID,
	})
	return respData, nil
}
