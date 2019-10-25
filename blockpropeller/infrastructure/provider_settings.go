package infrastructure

import (
	"context"
	"sync"
	"time"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/encryption"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	// ErrProviderSettingsNotFound is returned when a ProviderSettingsRepository does not find a provider settings to return.
	ErrProviderSettingsNotFound = errors.New("provider settings not found")
	// ErrProviderSettingsAlreadyExists is returned when a ProviderSettings creation is attempted with an existing ProviderSettingsID.
	ErrProviderSettingsAlreadyExists = errors.New("provider settings already exists")
)

// NilProviderSettingsID is an empty ProviderSettingsID.
var NilProviderSettingsID ProviderSettingsID

// ProviderSettingsID is a unique server identifier.
type ProviderSettingsID string

// NewProviderSettingsID returns a new unique ProviderSettingsID.
func NewProviderSettingsID() ProviderSettingsID {
	return ProviderSettingsID(uuid.NewV4().String())
}

// String satisfies the Stringer interface.
func (id ProviderSettingsID) String() string {
	return string(id)
}

// ProviderSettings hold access information to configured providers that the user
// has setup for his account. Only providers with valid settings can be used
// to provision new infrastructure.
type ProviderSettings struct {
	ID        ProviderSettingsID `json:"id" gorm:"varchar(36) not null"`
	AccountID account.ID         `json:"-" gorm:"type:varchar(36) not null references accounts(id)" `

	Label string `json:"label" gorm:"type:varchar(255) not null"`

	Type ProviderType `json:"type" gorm:"type:varchar(255) not null"`

	Credentials string `json:"-" gorm:"type:text not null"`

	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"-" gorm:"type:timestamp"`
}

// NewProviderSettings returns a new ProviderSettings instance.
func NewProviderSettings(accountID account.ID, label string, providerType ProviderType, credentials string) *ProviderSettings {
	return &ProviderSettings{
		ID:        NewProviderSettingsID(),
		AccountID: accountID,

		Label: label,

		Type:        providerType,
		Credentials: credentials,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BeforeSave encrypts the sensitive information before inserting them to the database.
func (ps *ProviderSettings) BeforeSave() error {
	encrypted, err := encryption.Encrypt([]byte(ps.Credentials))
	if err != nil {
		return errors.Wrap(err, "encrypt provider credentials")
	}

	ps.Credentials = string(encrypted)

	return nil
}

// AfterFind decrypts sensitive information before they can be used.
func (ps *ProviderSettings) AfterFind() error {
	decrypted, err := encryption.Decrypt([]byte(ps.Credentials))
	if err != nil {
		return errors.Wrap(err, "decrypt provider credentials")
	}

	ps.Credentials = string(decrypted)

	return nil
}

// BeforeDelete GORM hook to null out the credentials field.
func (ps *ProviderSettings) BeforeDelete(tx *gorm.DB) error {
	err := tx.Model(&ProviderSettings{}).
		Where("id = ?", ps.ID).
		Update("credentials", "[DELETED]").
		Error
	if err != nil {
		return errors.Wrap(err, "remove credentials")
	}

	return nil
}

// ProviderSettingsRepository defines an interface for storing and retrieving provisioning provider settings.
type ProviderSettingsRepository interface {
	// List all provider settings for a particular owner.
	List(ctx context.Context, ownerID account.ID) ([]*ProviderSettings, error)

	// Find a ProviderSettings given a ProviderSettingsID.
	Find(ctx context.Context, id ProviderSettingsID) (*ProviderSettings, error)

	// Create a new ProviderSettings.
	Create(ctx context.Context, providerSettings *ProviderSettings) error

	// Update an existing ProviderSettings.
	Update(ctx context.Context, providerSettings *ProviderSettings) error

	// Delete an existing ProviderSettings.
	Delete(ctx context.Context, providerSettings *ProviderSettings) error
}

// InMemoryProviderSettingsRepository holds the provider settings inside an in-memory map.
//
// ProviderSettings are not persisted on disk and won't survive program restarts.
type InMemoryProviderSettingsRepository struct {
	providerSettings sync.Map
}

// NewInMemoryProviderSettingsRepository returns a new InMemoryProviderSettingsRepository instance.
func NewInMemoryProviderSettingsRepository() *InMemoryProviderSettingsRepository {
	return &InMemoryProviderSettingsRepository{}
}

// List all provider settings for a particular owner.
func (repo *InMemoryProviderSettingsRepository) List(ctx context.Context, ownerID account.ID) ([]*ProviderSettings, error) {
	var settings []*ProviderSettings
	repo.providerSettings.Range(func(k, v interface{}) bool {
		setting := v.(*ProviderSettings)
		if setting.AccountID != ownerID {
			return true
		}

		settings = append(settings, setting)

		return true
	})

	return settings, nil
}

// Find a ProviderSettings given a ProviderSettingsID.
func (repo *InMemoryProviderSettingsRepository) Find(ctx context.Context, id ProviderSettingsID) (*ProviderSettings, error) {
	req, ok := repo.providerSettings.Load(id)
	if !ok {
		return nil, ErrProviderSettingsNotFound
	}

	return req.(*ProviderSettings), nil
}

// Create a new ProviderSettings.
func (repo *InMemoryProviderSettingsRepository) Create(ctx context.Context, providerSettings *ProviderSettings) error {
	_, loaded := repo.providerSettings.LoadOrStore(providerSettings.ID, providerSettings)
	if loaded {
		return ErrProviderSettingsAlreadyExists
	}

	return nil
}

// Update an existing ProviderSettings.
func (repo *InMemoryProviderSettingsRepository) Update(ctx context.Context, providerSettings *ProviderSettings) error {
	repo.providerSettings.Store(providerSettings.ID, providerSettings)

	return nil
}

// Delete an existing ProviderSettings.
func (repo *InMemoryProviderSettingsRepository) Delete(ctx context.Context, providerSettings *ProviderSettings) error {
	repo.providerSettings.Delete(providerSettings.ID)

	return nil
}
