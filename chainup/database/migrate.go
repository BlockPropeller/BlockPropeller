package database

import (
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"github.com/jinzhu/gorm"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
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
