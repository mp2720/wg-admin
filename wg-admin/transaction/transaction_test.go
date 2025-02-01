package transaction_test

import (
	"context"
	"errors"
	"mp2720/wg-admin/wg-admin/transaction"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockTransactionInitiator struct {
	tx         mockTransaction
	BeginError error
}

func (i *mockTransactionInitiator) Begin(ctx context.Context) (transaction.Tx, error) {
	return &i.tx, i.BeginError
}

func (i *mockTransactionInitiator) Reset() {
	*i = mockTransactionInitiator{}
}

type mockTransaction struct {
	CommitError, RollbackError   error
	CommitCalled, RollbackCalled int
}

func (tx *mockTransaction) Commit(ctx context.Context) error {
	tx.CommitCalled++
	return tx.CommitError
}

func (tx *mockTransaction) Rollback(ctx context.Context) error {
	tx.RollbackCalled++
	return tx.RollbackError
}

func Test_FromCtxWithNoTransaction(t *testing.T) {
	require.Nil(t, transaction.FromCtx(context.Background()))
}

func Test_Transaction(t *testing.T) {
	in := mockTransactionInitiator{}
	tm := transaction.NewManager(&in)

	// Test begin error.
	in.BeginError = errors.New("begin error")
	err := tm.InTransaction(context.Background(), func(ctx context.Context) error {
		t.FailNow()
		return nil
	})

	// commit
	{
		in.Reset()
		err = tm.InTransaction(context.Background(), func(ctx context.Context) error {
			// check tx is injected
			tx := transaction.FromCtx(ctx)
			require.IsType(t, &in.tx, tx)
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 1, in.tx.CommitCalled)
		require.Equal(t, 0, in.tx.RollbackCalled)

		// check commit error
		in.Reset()
		in.tx.CommitError = errors.New("Commit error")
		err = tm.InTransaction(context.Background(), func(ctx context.Context) error {
			return nil
		})
		require.ErrorIs(t, err, in.tx.CommitError)
		require.Equal(t, 1, in.tx.CommitCalled)
		require.Equal(t, 0, in.tx.RollbackCalled)

		// check nested
		in.Reset()
		err = tm.InTransaction(context.Background(), func(ctx context.Context) error {
			err := tm.InTransaction(ctx, func(ctx context.Context) error {
				return nil
			})
			require.Equal(t, 0, in.tx.CommitCalled)
			require.Equal(t, 0, in.tx.RollbackCalled)
			return err
		})
		require.NoError(t, err)
		require.Equal(t, 1, in.tx.CommitCalled)
		require.Equal(t, 0, in.tx.RollbackCalled)
	}

	// Test rollback
	{
		errorInTx := errors.New("Some error within the transaction")

		in.Reset()
		err = tm.InTransaction(context.Background(), func(ctx context.Context) error {
			return errorInTx
		})
		require.ErrorIs(t, err, errorInTx)
		require.Equal(t, 0, in.tx.CommitCalled)
		require.Equal(t, 1, in.tx.RollbackCalled)

		in.Reset()
		in.tx.RollbackError = errors.New("Rollback error")
		err = tm.InTransaction(context.Background(), func(ctx context.Context) error {
			return errorInTx
		})
		require.ErrorIs(t, err, in.tx.RollbackError)
		require.Equal(t, 0, in.tx.CommitCalled)
		require.Equal(t, 1, in.tx.RollbackCalled)

		// check nested
		in.Reset()
		err = tm.InTransaction(context.Background(), func(ctx context.Context) error {
			// check error is saved
			tm.InTransaction(ctx, func(ctx context.Context) error {
				return errorInTx
			})
			require.Equal(t, 0, in.tx.CommitCalled)
			require.Equal(t, 0, in.tx.RollbackCalled)
			return nil
		})
		require.ErrorIs(t, err, errorInTx)
		require.Equal(t, 0, in.tx.CommitCalled)
		require.Equal(t, 1, in.tx.RollbackCalled)
	}
}
