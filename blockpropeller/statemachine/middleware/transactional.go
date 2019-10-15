package middleware

import (
	"context"

	"blockpropeller.dev/blockpropeller/database/transaction"
	"blockpropeller.dev/blockpropeller/statemachine"
)

// Transactional wraps the execution of a Step with a database transaction.
type Transactional struct {
	txContext transaction.TxContext
}

// NewTransactional returns a new Transactional interface.
func NewTransactional(txContext transaction.TxContext) *Transactional {
	return &Transactional{txContext: txContext}
}

// Wrap implements the Middleware interface.
func (mw *Transactional) Wrap(step statemachine.Step) statemachine.Step {
	return statemachine.StepFn(
		func(ctx context.Context, res statemachine.StatefulResource) error {
			return mw.txContext.RunInTransaction(ctx, func(ctx context.Context) error {
				return step.Step(ctx, res)
			})
		},
	)
}
