package database

import (
	"context"

	"github.com/jinzhu/gorm"
)

type key int

var dbKey key

func newTxContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func txFromContext(ctx context.Context) (*gorm.DB, bool) {
	u, ok := ctx.Value(dbKey).(*gorm.DB)
	return u, ok
}

//DB is a wrapper structure over the *gorm.DB that adds the support for
// reusing the active transaction from the context.Context.
type DB struct {
	db *gorm.DB
}

// NewDB returns a new DB instance.
func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db.
			Set("gorm:association_autoupdate", false).
			Set("gorm:association_autocreate", false),
	}
}

// RunInTransaction starts a transaction and executes the provided callback in it.
//
//@TODO: Test rollback in case of panic.
func (db *DB) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, exists := txFromContext(ctx)
	if exists {
		panic("transaction already started")
	}

	tx := db.db.Begin()

	err := fn(newTxContext(ctx, tx))
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

// Model is a wrapper around the method for a gorm.DB method with the same name.
func (db *DB) Model(ctx context.Context, value interface{}) *gorm.DB {
	return db.getDB(ctx).Model(value)
}

// Close satisfies the Closer interface.
func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := txFromContext(ctx); ok {
		return tx
	}

	return db.db
}
