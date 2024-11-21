package cx

import (
	"assignment-pe/internal/log"
	"context"

	"github.com/jmoiron/sqlx"
)

type ContextKey string

const (
	ContextKeyLogger ContextKey = "logger"
	ContextKeyTx     ContextKey = "tx"
	ContextKeyUserID ContextKey = "userID"
)

func SetLogger(ctx context.Context, log log.AppLogger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, log)
}

func GetLogger(ctx context.Context) log.AppLogger {
	return ctx.Value(ContextKeyLogger).(log.AppLogger)
}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, ContextKeyTx, tx)
}

func GetTx(ctx context.Context) *sqlx.Tx {
	return ctx.Value(ContextKeyTx).(*sqlx.Tx)
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(ContextKeyUserID).(string)
}
