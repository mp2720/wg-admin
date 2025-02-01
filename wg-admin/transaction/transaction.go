package transaction

import (
	"context"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Something that begins a transaction.
type TransactionInitiator interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Manager struct {
	initiator TransactionInitiator
}

func (m *Manager) Begin(ctx context.Context) (Transaction, error) {
	return m.initiator.Begin(ctx)
}

type ctxKey struct{}

type transactionCtxWrapper struct {
	tx        Transaction
	lastError error
}

func extractFromCtx(ctx context.Context) *transactionCtxWrapper {
	tx := ctx.Value(ctxKey{})
	if tx == nil {
		return nil
	}
	return tx.(*transactionCtxWrapper)
}

func injectIntoCtx(ctx context.Context, tx *transactionCtxWrapper) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

// Run f inside a transaction injected into the context.
//
// Note that nested calls do not begin a new transaction - all are executed inside the same one.
// If in any nested call f fails, then the whole transaction is rolled back, otherwise commited.
//
// Returns the error from last failed f call or the rollback/commit error.
func (m *Manager) InTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	txInParentCtx := extractFromCtx(ctx)
	if txInParentCtx != nil {
		// there's already ongoing transaction in this context - just save the error and return it
		txInParentCtx.lastError = f(ctx)
		return txInParentCtx.lastError
	}

	// it's the first call in context, there's no ongoing transaction

	newTx, err := m.Begin(ctx)
	if err != nil {
		return err
	}

	newTxWrapper := transactionCtxWrapper{
		tx:        newTx,
		lastError: nil,
	}

	newTxWrapper.lastError = f(injectIntoCtx(ctx, &newTxWrapper))

	// do commit/rollback

	if newTxWrapper.lastError != nil {
		return newTx.Rollback(ctx)
	} else {
		return newTx.Commit(ctx)
	}
}
