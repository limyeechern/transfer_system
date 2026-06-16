package router

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"transfer_system/biz/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

const logIDHeader = "X-Log-ID"

func RootMiddleware() []app.HandlerFunc {
	return []app.HandlerFunc{LogIDMiddleware()}
}

func LogIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logID := string(c.GetHeader(logIDHeader))
		if logID == "" {
			logID = newLogID()
		}
		c.Header(logIDHeader, logID)
		c.Next(context.WithValue(ctx, consts.LogID, logID))
	}
}

func newLogID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b[:])
}

func AccountsMiddleware() []app.HandlerFunc {
	return nil
}

func TransactionsMiddleware() []app.HandlerFunc {
	return nil
}
