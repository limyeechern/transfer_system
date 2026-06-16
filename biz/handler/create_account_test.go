package handler_test

import (
	"context"
	"errors"
	"fmt"
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

type fakeCreateAccountService struct {
	account *model.Account
	err     error
	req     *model.NewAccount
}

func (s *fakeCreateAccountService) Create(ctx context.Context, req *model.NewAccount) (*model.Account, error) {
	s.req = req
	if s.err != nil {
		return nil, s.err
	}
	return s.account, nil
}

func (s *fakeCreateAccountService) Validate(ctx context.Context, req *model.NewAccount) error {
	return s.err
}

func TestCreateAccountEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		serviceErr error
		wantStatus int
		wantBody   string
		wantReq    *model.NewAccount
	}{
		{
			name:       "success",
			body:       `{"account_id":123,"initial_balance":"100.23344"}`,
			wantStatus: consts.StatusOK,
			wantBody:   `{}`,
			wantReq: &model.NewAccount{
				AccountID:      123,
				InitialBalance: "100.23344",
			},
		},
		{
			name:       "duplicate account id",
			body:       `{"account_id":123,"initial_balance":"100.23344"}`,
			serviceErr: apperror.ErrAccountIdAlreadyExists,
			wantStatus: consts.StatusBadRequest,
			wantBody:   errorJSON("ACCOUNT_ID_ALREADY_EXISTS", apperror.ErrAccountIdAlreadyExists.Error()),
			wantReq: &model.NewAccount{
				AccountID:      123,
				InitialBalance: "100.23344",
			},
		},
		{
			name:       "invalid initial balance",
			body:       `{"account_id":124,"initial_balance":"100.233445"}`,
			serviceErr: apperror.ErrInvalidAmount,
			wantStatus: consts.StatusBadRequest,
			wantBody:   errorJSON("INVALID_AMOUNT", apperror.ErrInvalidAmount.Error()),
			wantReq: &model.NewAccount{
				AccountID:      124,
				InitialBalance: "100.233445",
			},
		},
		{
			name:       "invalid json",
			body:       `{"account_id":125,"initial_balance":`,
			wantStatus: consts.StatusBadRequest,
			wantBody:   errorJSON("INVALID_REQUEST", apperror.ErrInvalidRequest.Error()),
		},
		{
			name:       "internal error",
			body:       `{"account_id":126,"initial_balance":"100.23344"}`,
			serviceErr: apperror.ErrInternalError,
			wantStatus: consts.StatusInternalServerError,
			wantBody:   errorJSON("INTERNAL_ERROR", apperror.ErrInternalError.Error()),
			wantReq: &model.NewAccount{
				AccountID:      126,
				InitialBalance: "100.23344",
			},
		},
		{
			name:       "unknown error",
			body:       `{"account_id":127,"initial_balance":"100.23344"}`,
			serviceErr: errors.New("database connection failed"),
			wantStatus: consts.StatusInternalServerError,
			wantBody:   errorJSON("INTERNAL_ERROR", apperror.ErrInternalError.Error()),
			wantReq: &model.NewAccount{
				AccountID:      127,
				InitialBalance: "100.23344",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createAccountService := &fakeCreateAccountService{
				account: &model.Account{AccountID: 123, Balance: "10023344"},
				err:     tt.serviceErr,
			}
			h := server.Default()
			router.Register(h, handler.NewApp(handler.Dependencies{
				CreateAccountService: createAccountService,
			}))

			resp := ut.PerformRequest(
				h.Engine,
				consts.MethodPost,
				"/accounts",
				&ut.Body{Body: strings.NewReader(tt.body), Len: len(tt.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)

			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%q", resp.Code, tt.wantStatus, resp.Body.String())
			}
			if tt.wantBody != "" && strings.TrimSpace(resp.Body.String()) != tt.wantBody {
				t.Fatalf("body = %q, want %q", strings.TrimSpace(resp.Body.String()), tt.wantBody)
			}
			if tt.wantReq == nil {
				return
			}
			if createAccountService.req == nil {
				t.Fatalf("service did not receive request")
			}
			if *createAccountService.req != *tt.wantReq {
				t.Fatalf("service req = %+v, want %+v", createAccountService.req, tt.wantReq)
			}
		})
	}
}

func errorJSON(code string, message string) string {
	return fmt.Sprintf(`{"code":"%s","message":"%s"}`, code, message)
}
