package logs

import (
	"context"

	"transfer_system/biz/consts"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

func LogWithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		string(consts.LogID): ctx.Value(consts.LogID),
	})
}

func CtxInfo(ctx context.Context, msg string, fields ...Fields) {
	withFields(ctx, fields...).Info(msg)
}

func CtxError(ctx context.Context, msg string, err error, fields ...Fields) {
	entry := withFields(ctx, fields...)
	if err != nil {
		entry = entry.WithError(err)
	}
	entry.Error(msg)
}

func CtxWarn(ctx context.Context, msg string, fields ...Fields) {
	withFields(ctx, fields...).Warn(msg)
}

func withFields(ctx context.Context, fields ...Fields) *logrus.Entry {
	entry := LogWithContext(ctx)
	for _, fieldSet := range fields {
		logrusFields := logrus.Fields{}
		for key, value := range fieldSet {
			logrusFields[key] = value
		}
		entry = entry.WithFields(logrusFields)
	}
	return entry
}
