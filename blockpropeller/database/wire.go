package database

import (
	"fmt"

	"blockpropeller.dev/blockpropeller/database/transaction"
	"blockpropeller.dev/lib/log"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import driver for the SQL dialect.
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // Import driver for the SQL dialect.
	"github.com/pkg/errors"
)

// Set keeps a set of all database level dependencies.
var Set = wire.NewSet(
	ProvideDB,
	wire.Bind(new(transaction.TxContext), new(*DB)),
)

// ProvideDB initializes and returns a new DB instance.
//
// Along with connecting to the database, ProvideDB also runs
// the auto migration process to create any potentially missing tables or columns.
//
// A cleanup function is returned with the database client so cleanup can be done efficiently.
func ProvideDB(cfg *Config, logCfg *log.Config) (db *DB, closeFn func(), err error) {
	switch cfg.Dialect {
	case "sqlite3":
		db, err = provideSqliteDB(cfg)
	case "postgres":
		db, err = providePostgresDB(cfg)
	default:
		err = errors.Errorf("unsupported database dialect: %s", cfg.Dialect)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "provide DB")
	}

	if logCfg.Level == log.LevelDebug {
		db.db.LogMode(true)
	}

	err = migrate(db.db)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed auto migrating models")
	}

	return db, func() {
		log.Info("closing database client...")
		log.Closer(db)
	}, nil
}

func provideSqliteDB(cfg *Config) (*DB, error) {
	db, err := gorm.Open("sqlite3", cfg.File)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlite3 database")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	return NewDB(db), nil
}

func providePostgresDB(cfg *Config) (*DB, error) {
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Database))
	if err != nil {
		return nil, errors.Wrap(err, "open postgres database")
	}

	return NewDB(db), nil
}
