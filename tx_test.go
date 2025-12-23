package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/momoli-dev/mogo/database"
)

func TestNoTx(t *testing.T) {
	ctx := context.Background()

	inTx := database.InTx(ctx)
	require.False(t, inTx)

	tx := database.Tx(ctx)
	require.Nil(t, tx)
}

func TestBeginTx_OK(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		newCtx, err := conn.BeginTx(ctx)
		require.NoError(t, err)

		inTx := database.InTx(newCtx)
		require.True(t, inTx)

		tx := database.Tx(newCtx)
		require.NotNil(t, tx)
	})
}

func TestBeginTx_AlreadyInTx(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		newCtx, err := conn.BeginTx(ctx)
		require.NoError(t, err)

		sameCtx, err := conn.BeginTx(newCtx)
		require.NoError(t, err)
		require.Equal(t, newCtx, sameCtx)

		inTx := database.InTx(sameCtx)
		require.True(t, inTx)
	})
}

func TestCommitTx_OK(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		newCtx, err := conn.BeginTx(ctx)
		require.NoError(t, err)

		err = database.CommitTx(newCtx)
		require.NoError(t, err)
	})
}

func TestCommitTx_NoTx(t *testing.T) {
	RunWithConn(t, func(_ *database.Conn) {
		ctx := context.Background()
		err := database.CommitTx(ctx)
		require.NoError(t, err)
	})
}

func TestRollbackTx_OK(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		newCtx, err := conn.BeginTx(ctx)
		require.NoError(t, err)

		err = database.RollbackTx(newCtx)
		require.NoError(t, err)
	})
}

func TestRollbackTx_NoTx(t *testing.T) {
	RunWithConn(t, func(_ *database.Conn) {
		ctx := context.Background()
		err := database.RollbackTx(ctx)
		require.NoError(t, err)
	})
}

func TestWithTx_OK(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		err := conn.WithTx(ctx, func(txCtx context.Context) error {
			inTx := database.InTx(txCtx)
			require.True(t, inTx)

			tx := database.Tx(txCtx)
			require.NotNil(t, tx)

			return nil
		})
		require.NoError(t, err)
	})
}

func TestWithTx_NestedTx(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		err := conn.WithTx(ctx, func(txCtx context.Context) error {
			inTx := database.InTx(txCtx)
			require.True(t, inTx)

			tx := database.Tx(txCtx)
			require.NotNil(t, tx)

			return conn.WithTx(txCtx, func(nestedTxCtx context.Context) error {
				nestedInTx := database.InTx(nestedTxCtx)
				require.True(t, nestedInTx)

				nestedTx := database.Tx(nestedTxCtx)
				require.NotNil(t, nestedTx)

				require.Equal(t, tx, nestedTx)
				require.Equal(t, txCtx, nestedTxCtx)

				return nil
			})
		})
		require.NoError(t, err)
	})
}

func TestWithTx_ErrorInFunc(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		err := conn.WithTx(ctx, func(_ context.Context) error {
			return assert.AnError
		})
		require.Error(t, err)
	})
}

func TestWithTxHelper_OK(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		result, err := database.WithTx(ctx, conn, func(txCtx context.Context) (string, error) {
			inTx := database.InTx(txCtx)
			require.True(t, inTx)

			tx := database.Tx(txCtx)
			require.NotNil(t, tx)

			return "success", nil
		})
		require.NoError(t, err)
		require.Equal(t, "success", result)
	})
}

func TestWithTxHelper_ErrorInFunc(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		result, err := database.WithTx(ctx, conn, func(_ context.Context) (string, error) {
			return "", assert.AnError
		})
		require.Error(t, err)
		require.Empty(t, result)
	})
}

func TestWithTxHelper_Nested(t *testing.T) {
	RunWithConn(t, func(conn *database.Conn) {
		ctx := context.Background()
		res, err := database.WithTx(ctx, conn, func(txCtx context.Context) (string, error) {
			inTx := database.InTx(txCtx)
			require.True(t, inTx)

			tx := database.Tx(txCtx)
			require.NotNil(t, tx)

			return database.WithTx(txCtx, conn, func(nestedTxCtx context.Context) (string, error) {
				nestedInTx := database.InTx(nestedTxCtx)
				require.True(t, nestedInTx)

				nestedTx := database.Tx(nestedTxCtx)
				require.NotNil(t, nestedTx)

				require.Equal(t, tx, nestedTx)
				require.Equal(t, txCtx, nestedTxCtx)

				return "nested success", nil
			})
		})
		require.Equal(t, "nested success", res)
		require.NoError(t, err)
	})
}
