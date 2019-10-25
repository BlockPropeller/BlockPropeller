package database

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"github.com/jinzhu/gorm"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&account.Account{},
		&infrastructure.ProviderSettings{},
		&infrastructure.Server{},
		&infrastructure.Deployment{},
		&provision.Job{},
	).Error
	if err != nil {
		return err
	}

	return nil
}
