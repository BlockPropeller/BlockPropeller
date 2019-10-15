package database

import (
	"context"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/infrastructure"
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

// List all provider settings for a particular owner.
func (repo *ProviderSettingsRepository) List(ctx context.Context, ownerID account.ID) ([]*infrastructure.ProviderSettings, error) {
	var settings []*infrastructure.ProviderSettings
	err := repo.db.Model(ctx, &settings).
		Where("account_id = ?", ownerID).
		Find(&settings).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "list provider settings")
	}

	return settings, nil
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
