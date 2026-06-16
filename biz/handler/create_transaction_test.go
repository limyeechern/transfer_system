package handler_test

import (
	"context"
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

type fakeCreateTransactionService struct {
	transaction *model.Transaction
	err         error
	req         *model.Transaction
}

func (s *fakeCreateTransactionService) Create(ctx context.Context, req *model.Transaction) (*model.Transaction, error) {
	s.req = req
	if s.err != nil {
		return nil, s.err
	}
	return s.transaction, nil
}

func (s *fakeCreateTransactionService) Validate(ctx context.Context, req *model.Transaction) error {
	return s.err
}

func TestCreateTransactionEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		serviceErr error
		wantStatus int
		wantBody   string
		wantReq    *model.Transaction
	}{
		{
			name:       "success",
			body:       `{"source_account_id":123,"destination_account_id":456,"amount":"100.12345"}`,
			wantStatus: consts.StatusOK,
			wantBody:   `{}`,
			wantReq: &model.Transaction{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "100.12345",
			},
		},
		{
			name:       "invalid amount",
			body:       `{"source_account_id":123,"destination_account_id":456,"amount":"0"}`,
			serviceErr: apperror.ErrInvalidAmount,
			wantStatus: consts.StatusBadRequest,
			wantBody:   apperror.ErrInvalidAmount.Error(),
			wantReq: &model.Transaction{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "0",
			},
		},
		{
			name:       "account not found",
			body:       `{"source_account_id":999,"destination_account_id":456,"amount":"100.12345"}`,
			serviceErr: apperror.ErrAccountNotFound,
			wantStatus: consts.StatusNotFound,
			wantBody:   apperror.ErrAccountNotFound.Error(),
			wantReq: &model.Transaction{
				SourceAccountID:      999,
				DestinationAccountID: 456,
				Amount:               "100.12345",
			},
		},
		{
			name:       "insufficient balance",
			body:       `{"source_account_id":123,"destination_account_id":456,"amount":"100.12345"}`,
			serviceErr: apperror.ErrInsufficientBalance,
			wantStatus: consts.StatusBadRequest,
			wantBody:   apperror.ErrInsufficientBalance.Error(),
			wantReq: &model.Transaction{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "100.12345",
			},
		},
		{
			name:       "invalid transaction",
			body:       `{"source_account_id":123,"destination_account_id":123,"amount":"100.12345"}`,
			serviceErr: apperror.ErrInvalidTransaction,
			wantStatus: consts.StatusBadRequest,
			wantBody:   apperror.ErrInvalidTransaction.Error(),
			wantReq: &model.Transaction{
				SourceAccountID:      123,
				DestinationAccountID: 123,
				Amount:               "100.12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTransactionService := &fakeCreateTransactionService{
				transaction: &model.Transaction{SourceAccountID: 123, DestinationAccountID: 456, Amount: "100.12345"},
				err:         tt.serviceErr,
			}
			h := server.Default()
			router.Register(h, handler.NewApp(handler.Dependencies{
				CreateTransactionService: createTransactionService,
			}))

			resp := ut.PerformRequest(
				h.Engine,
				consts.MethodPost,
				"/transactions",
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
			if createTransactionService.req == nil {
				t.Fatalf("service did not receive request")
			}
			if *createTransactionService.req != *tt.wantReq {
				t.Fatalf("service req = %+v, want %+v", createTransactionService.req, tt.wantReq)
			}
		})
	}
}
