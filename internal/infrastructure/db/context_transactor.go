package db

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type ctxTransactor struct{}

func addCtxTransact(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, ctxTransactor{}, tx)
}

func extractCtxTransact(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(ctxTransactor{}).(*sqlx.Tx); ok {
		return tx
	}

	return nil
}
