package handler_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"transfer_system/biz/apperror"
	"transfer_system/biz/handler"
	"transfer_system/biz/model"
	"transfer_system/biz/router"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type fakeGetAccountService struct {
	account *model.Account
	err     error
	req     *model.GetAccount
}

func (s *fakeGetAccountService) Read(ctx context.Context, req *model.GetAccount) (*model.Account, error) {
	s.req = req
	if s.err != nil {
		return nil, s.err
	}
	return s.account, nil
}

func (s *fakeGetAccountService) Validate(ctx context.Context, req *model.GetAccount) error {
	return s.err
}

func TestGetAccountEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		serviceErr error
		wantStatus int
		wantBody   string
		wantReq    *model.GetAccount
	}{
		{
			name:       "success",
			path:       "/accounts/123",
			wantStatus: consts.StatusOK,
			wantBody:   `{"account_id":123,"balance":"100.23344"}`,
			wantReq:    &model.GetAccount{AccountID: 123},
		},
		{
			name:       "account not found",
			path:       "/accounts/123",
			serviceErr: apperror.ErrAccountNotFound,
			wantStatus: consts.StatusNotFound,
			wantBody:   apperror.ErrAccountNotFound.Error(),
			wantReq:    &model.GetAccount{AccountID: 123},
		},
		{
			name:       "invalid account id",
			path:       "/accounts/not-an-id",
			wantStatus: consts.StatusBadRequest,
			wantBody:   apperror.ErrInvalidAccount.Error(),
		},
		{
			name:       "internal error",
			path:       "/accounts/123",
			serviceErr: apperror.ErrInternalError,
			wantStatus: consts.StatusInternalServerError,
			wantBody:   apperror.ErrInternalError.Error(),
			wantReq:    &model.GetAccount{AccountID: 123},
		},
		{
			name:       "unknown error",
			path:       "/accounts/123",
			serviceErr: errors.New("database connection failed"),
			wantStatus: consts.StatusInternalServerError,
			wantBody:   apperror.ErrInternalError.Error(),
			wantReq:    &model.GetAccount{AccountID: 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAccountService := &fakeGetAccountService{
				account: &model.Account{AccountID: 123, Balance: "100.23344"},
				err:     tt.serviceErr,
			}
			h := server.Default()
			router.Register(h, handler.NewApp(handler.Dependencies{
				GetAccountService: getAccountService,
			}))

			resp := ut.PerformRequest(h.Engine, consts.MethodGet, tt.path, nil)

			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%q", resp.Code, tt.wantStatus, resp.Body.String())
			}
			if tt.wantBody != "" && strings.TrimSpace(resp.Body.String()) != tt.wantBody {
				t.Fatalf("body = %q, want %q", strings.TrimSpace(resp.Body.String()), tt.wantBody)
			}
			if tt.wantReq == nil {
				return
			}
			if getAccountService.req == nil {
				t.Fatalf("service did not receive request")
			}
			if *getAccountService.req != *tt.wantReq {
				t.Fatalf("service req = %+v, want %+v", getAccountService.req, tt.wantReq)
			}
		})
	}
}
