package handler

import (
	"errors"

	"transfer_system/biz/apperror"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func writeError(c *app.RequestContext, err error) {
	c.String(statusCodeForError(err), messageForError(err))
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
