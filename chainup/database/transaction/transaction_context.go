package transaction

import "context"

// TxContext provides a way to execute code that is scoped to a transaction.
//
// This allows us to rollback any changes made in case of failure.
type TxContext interface {
	RunInTransaction(context.Context, func(context.Context) error) error
}

// InMemoryTxContext is an in-memory implementation of a database.InMemoryTxContext.
//
// This implementation panics on errors because we cannot easily handle transactional
// repositories in memory.
type InMemoryTxContext struct {
}

// NewInMemoryTransactionContext returns a new InMemoryTxContext instance.
func NewInMemoryTransactionContext() *InMemoryTxContext {
	return &InMemoryTxContext{}
}

// RunInTransaction satisfies the TxContext interface.
func (InMemoryTxContext) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := fn(ctx)
	if err != nil {
		panic("failed executing callback in a transaction")
	}

	return nil
}
