package account

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

var (
	// ErrAccountNotFound is returned when an account.Repository does not find an account to return.
	ErrAccountNotFound = errors.New("account not found")
)

// Repository defines an interface for storing and retrieving accounts.
type Repository interface {
	// FindByID returns an Account with the given ID.
	FindByID(ctx context.Context, id ID) (*Account, error)
	// FindByEmail returns an Account with the provided email.
	FindByEmail(ctx context.Context, email Email) (*Account, error)
	// List all accounts.
	List(ctx context.Context) ([]*Account, error)

	// Create inserts a new Account into the repository.
	Create(ctx context.Context, acc *Account) error
	// Update updates an existing account in the repository.
	Update(ctx context.Context, acc *Account) error
}

// InMemoryRepository is an in-memory implementation of an account.Repository.
//
// Accounts are not persisted on disk and won't survive program restarts.
type InMemoryRepository struct {
	accountsByID    sync.Map
	accountsByEmail sync.Map
}

// NewInMemoryRepository returns a new InMemoryRepository instance.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

// FindByID returns an Account with the given ID.
func (repo *InMemoryRepository) FindByID(ctx context.Context, id ID) (*Account, error) {
	req, ok := repo.accountsByID.Load(id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	return req.(*Account), nil
}

// FindByEmail returns an Account with the provided email.
func (repo *InMemoryRepository) FindByEmail(ctx context.Context, email Email) (*Account, error) {
	req, ok := repo.accountsByEmail.Load(email)
	if !ok {
		return nil, ErrAccountNotFound
	}

	return req.(*Account), nil
}

// List all accounts.
func (repo *InMemoryRepository) List(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	repo.accountsByID.Range(func(k, v interface{}) bool {
		accounts = append(accounts, v.(*Account))

		return true
	})

	return accounts, nil
}

// Create inserts a new Account into the repository.
func (repo *InMemoryRepository) Create(ctx context.Context, acc *Account) error {
	return repo.set(acc)
}

// Update updates an existing account in the repository.
func (repo *InMemoryRepository) Update(ctx context.Context, acc *Account) error {
	return repo.set(acc)
}

func (repo *InMemoryRepository) set(acc *Account) error {
	repo.accountsByID.Store(acc.ID, acc)
	repo.accountsByEmail.Store(acc.Email, acc)

	return nil
}
