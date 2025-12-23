package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Tx retrieves transaction from context.
//
//nolint:ireturn // this is aimed at pgx
func Tx(ctx context.Context) pgx.Tx {
	return txFromCtx(ctx)
}

// InTx checks if the context has an active transaction.
func InTx(ctx context.Context) bool {
	return Tx(ctx) != nil
}

// BeginTx starts a new transaction and returns a new context containing the transaction.
func (conn *Conn) BeginTx(ctx context.Context) (context.Context, error) {
	if tx := Tx(ctx); tx != nil {
		return ctx, nil
	}

	tx, err := conn.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	ctx = ctxWithTx(ctx, tx)

	return ctx, nil
}

// CommitTx commits the transaction stored in the context, if any.
func CommitTx(ctx context.Context) error {
	tx := Tx(ctx)
	if tx == nil {
		return nil
	}

	return tx.Commit(ctx)
}

// RollbackTx rolls back the transaction stored in the context, if any.
func RollbackTx(ctx context.Context) error {
	tx := txFromCtx(ctx)
	if tx == nil {
		return nil
	}

	return tx.Rollback(ctx)
}

// WithTx executes the given function within a transaction context, starting a new transaction if one does not already exist.
func (conn *Conn) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if InTx(ctx) {
		return fn(ctx)
	}

	newCtx, err := conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := RollbackTx(newCtx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			//nolint:forbidigo // idk what to do here yet
			fmt.Println("failed to rollback transaction:", err)
		}
	}()

	if err := fn(newCtx); err != nil {
		return err
	}

	if err := CommitTx(newCtx); err != nil {
		return err
	}

	return nil
}

// WithTx is a generic helper version of conn.WithTx that can be used to return a value from the function as well.
func WithTx[T any](ctx context.Context, conn *Conn, fn func(ctx context.Context) (T, error)) (T, error) {
	var innerRes T
	err := conn.WithTx(ctx, func(ctx context.Context) error {
		res, err := fn(ctx)
		innerRes = res
		return err
	})
	return innerRes, err
}
