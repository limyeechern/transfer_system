package handler

import (
	"errors"

	"transfer_system/biz/apperror"
	"transfer_system/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func writeError(c *app.RequestContext, err error) {
	c.JSON(statusCodeForError(err), model.ErrorResponse{
		Code:    codeForError(err),
		Message: messageForError(err),
	})
}

func statusCodeForError(err error) int {
	switch {
	case errors.Is(err, apperror.ErrInvalidRequest),
		errors.Is(err, apperror.ErrInvalidAccount),
		errors.Is(err, apperror.ErrInvalidTransaction),
		errors.Is(err, apperror.ErrInvalidAmount),
		errors.Is(err, apperror.ErrAccountIdAlreadyExists),
		errors.Is(err, apperror.ErrInsufficientBalance):
		return consts.StatusBadRequest
	case errors.Is(err, apperror.ErrAccountNotFound):
		return consts.StatusNotFound
	default:
		return consts.StatusInternalServerError
	}
}

func codeForError(err error) string {
	switch {
	case errors.Is(err, apperror.ErrInvalidRequest):
		return "INVALID_REQUEST"
	case errors.Is(err, apperror.ErrInvalidAccount):
		return "INVALID_ACCOUNT"
	case errors.Is(err, apperror.ErrInvalidTransaction):
		return "INVALID_TRANSACTION"
	case errors.Is(err, apperror.ErrInvalidAmount):
		return "INVALID_AMOUNT"
	case errors.Is(err, apperror.ErrAccountIdAlreadyExists):
		return "ACCOUNT_ID_ALREADY_EXISTS"
	case errors.Is(err, apperror.ErrInsufficientBalance):
		return "INSUFFICIENT_BALANCE"
	case errors.Is(err, apperror.ErrAccountNotFound):
		return "ACCOUNT_NOT_FOUND"
	case errors.Is(err, apperror.ErrInternalError):
		return "INTERNAL_ERROR"
	default:
		return "INTERNAL_ERROR"
	}
}

func messageForError(err error) string {
	switch {
	case errors.Is(err, apperror.ErrInvalidRequest),
		errors.Is(err, apperror.ErrInvalidAccount),
		errors.Is(err, apperror.ErrInvalidTransaction),
		errors.Is(err, apperror.ErrInvalidAmount),
		errors.Is(err, apperror.ErrAccountIdAlreadyExists),
		errors.Is(err, apperror.ErrInsufficientBalance),
		errors.Is(err, apperror.ErrAccountNotFound),
		errors.Is(err, apperror.ErrInternalError):
		return err.Error()
	default:
		return apperror.ErrInternalError.Error()
	}
}
