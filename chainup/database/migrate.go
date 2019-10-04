package database

import (
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&account.Account{},
		&provision.Job{},
		&infrastructure.ProviderSettings{},
		&infrastructure.Server{},
		&infrastructure.Deployment{},
	).Error
	if err != nil {
		return err
	}

	err = db.Exec(`
CREATE UNIQUE INDEX IF NOT EXISTS account_email_idx ON accounts (email)
`).Error
	if err != nil {
		return errors.Wrap(err, "create indexes")
	}

	return nil
}
