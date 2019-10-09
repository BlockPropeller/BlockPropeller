package database

import (
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"github.com/jinzhu/gorm"
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

	return nil
}
