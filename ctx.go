package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// ctxKeyTx is the context key for stored transactions.
type ctxKeyTx struct{}

// txFromCtx retrieves a transaction from the context, if present. Returns nil if no transaction is found.
//
//nolint:ireturn // this is aimed at pgx
func txFromCtx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(ctxKeyTx{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}

// ctxWithTx returns a new context with the provided transaction stored in it.
func ctxWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}
