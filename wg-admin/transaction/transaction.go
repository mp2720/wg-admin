package transaction

import (
	"context"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Something that begins a transaction.
type Initiator interface {
	Begin(ctx context.Context) (Tx, error)
}

type Manager struct {
	initiator Initiator
}

func NewManager(initiator Initiator) Manager {
	return Manager{initiator}
}

func (m *Manager) Begin(ctx context.Context) (Tx, error) {
	return m.initiator.Begin(ctx)
}

type ctxKey struct{}

type ctxWrapper struct {
	tx        Tx
	lastError error
}

func extractWrapperFromCtx(ctx context.Context) *ctxWrapper {
	w := ctx.Value(ctxKey{})
	if w == nil {
		return nil
	}
	return w.(*ctxWrapper)
}

func FromCtx(ctx context.Context) Tx {
	txInCtx := extractWrapperFromCtx(ctx)
	if txInCtx == nil {
		return nil
	}
	return txInCtx.tx
}

func injectWrapperIntoCtx(ctx context.Context, w *ctxWrapper) context.Context {
	println("inject", w)
	return context.WithValue(ctx, ctxKey{}, w)
}

// Run f inside a transaction injected into the context.
//
// Note that nested calls do not begin a new transaction - all are executed inside the same one.
// If in any nested call f fails, then the whole transaction is rolled back, otherwise commited.
//
// Returns the error from last failed f call or the rollback/commit error.
func (m *Manager) InTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	txInParentCtx := extractWrapperFromCtx(ctx)
	println("extract", txInParentCtx)
	if txInParentCtx != nil {
		// there's already ongoing transaction in this context - just save the error and return it
		txInParentCtx.lastError = f(ctx)
		println("extract err", txInParentCtx.lastError)
		return txInParentCtx.lastError
	}

	// it's the first call in context, there's no ongoing transaction

	newTx, err := m.Begin(ctx)
	if err != nil {
		return err
	}

	newTxWrapper := ctxWrapper{
		tx:        newTx,
		lastError: nil,
	}

	err = f(injectWrapperIntoCtx(ctx, &newTxWrapper))
	if err != nil {
		newTxWrapper.lastError = err
	}

	// do commit/rollback

	println(&newTxWrapper, newTxWrapper.lastError)

	if newTxWrapper.lastError != nil {
		rollbackErr := newTx.Rollback(ctx)
		if rollbackErr != nil {
			return rollbackErr
		}

		return newTxWrapper.lastError
	} else {
		return newTx.Commit(ctx)
	}
}
