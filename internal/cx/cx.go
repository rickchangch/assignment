package cx

import (
	"assignment-pe/internal/log"
	"context"

	"github.com/jmoiron/sqlx"
)

type ContextKey string

const (
	ContextKeyLogger ContextKey = "logger"
	txKey            ContextKey = "tx"
)

func SetLogger(ctx context.Context, log log.AppLogger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, log)
}

func GetLogger(ctx context.Context) log.AppLogger {
	return ctx.Value(ContextKeyLogger).(log.AppLogger)
}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) *sqlx.Tx {
	return ctx.Value(txKey).(*sqlx.Tx)
}
