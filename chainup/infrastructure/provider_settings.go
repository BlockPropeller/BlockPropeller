package infrastructure

import (
	"context"
	"sync"
	"time"

	"chainup.dev/chainup/account"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	// ErrProviderSettingsNotFound is returned when a ProviderSettingsRepository does not find a provider settings to return.
	ErrProviderSettingsNotFound = errors.New("provider settings not found")
	// ErrProviderSettingsAlreadyExists is returned when a ProviderSettings creation is attempted with an existing ProviderSettingsID.
	ErrProviderSettingsAlreadyExists = errors.New("provider settings already exists")
)

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
	ID        ProviderSettingsID `json:"id"`
	AccountID account.ID         `json:"-" sql:"type:varchar(255) NOT NULL REFERENCES accounts(id)" `

	Type ProviderType `json:"type"`

	Credentials string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewProviderSettings returns a new ProviderSettings instance.
func NewProviderSettings(accountID account.ID, providerType ProviderType, credentials string) *ProviderSettings {
	return &ProviderSettings{
		ID:        NewProviderSettingsID(),
		AccountID: accountID,

		Type:        providerType,
		Credentials: credentials,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ProviderSettingsRepository defines an interface for storing and retrieving provisioning provider settings.
type ProviderSettingsRepository interface {
	// Find a ProviderSettings given a ProviderSettingsID.
	Find(ctx context.Context, id ProviderSettingsID) (*ProviderSettings, error)

	// Create a new ProviderSettings.
	Create(ctx context.Context, providerSettings *ProviderSettings) error

	// Update an existing ProviderSettings.
	Update(ctx context.Context, providerSettings *ProviderSettings) error
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
