package database

import (
	"context"

	"chainup.dev/chainup/infrastructure"
	"github.com/pkg/errors"
)

// ProviderSettingsRepository is a databased backed implementation of a infrastructure.ProviderSettingsRepository.
type ProviderSettingsRepository struct {
	db *DB
}

// NewProviderSettingsRepository returns a new ProviderSettingsRepository instance.
func NewProviderSettingsRepository(db *DB) *ProviderSettingsRepository {
	return &ProviderSettingsRepository{db: db}
}

// Find a ProviderSettings given a ProviderSettingsID.
func (repo *ProviderSettingsRepository) Find(ctx context.Context, id infrastructure.ProviderSettingsID) (*infrastructure.ProviderSettings, error) {
	var settings infrastructure.ProviderSettings
	err := repo.db.Model(ctx, &settings).Where("id = ?", id).First(&settings).Error
	if err != nil {
		return nil, errors.Wrap(err, "find provider settings by ID")
	}

	return &settings, nil
}

// Create a new ProviderSettings.
func (repo *ProviderSettingsRepository) Create(ctx context.Context, providerSettings *infrastructure.ProviderSettings) error {
	err := repo.db.Model(ctx, providerSettings).Create(providerSettings).Error
	if err != nil {
		return errors.Wrap(err, "create provider settings")
	}

	return nil
}

// Update an existing ProviderSettings.
func (repo *ProviderSettingsRepository) Update(ctx context.Context, providerSettings *infrastructure.ProviderSettings) error {
	err := repo.db.Model(ctx, providerSettings).Save(providerSettings).Error
	if err != nil {
		return errors.Wrap(err, "update provider settings")
	}

	return nil
}
