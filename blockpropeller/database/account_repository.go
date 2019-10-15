package database

import (
	"context"

	"blockpropeller.dev/blockpropeller/account"
	"github.com/pkg/errors"
)

// AccountRepository is a databased backed implementation of a account.Repository.
type AccountRepository struct {
	db *DB
}

// NewAccountRepository returns a new AccountRepository instance.
func NewAccountRepository(db *DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// FindByID returns an Account with the given ID.
func (repo *AccountRepository) FindByID(ctx context.Context, id account.ID) (*account.Account, error) {
	var acc account.Account
	err := repo.db.Model(ctx, &acc).Where("id = ?", id).First(&acc).Error
	if err != nil {
		return nil, errors.Wrap(err, "find account by ID")
	}

	return &acc, nil
}

// FindByEmail returns an Account with the provided email.
func (repo *AccountRepository) FindByEmail(ctx context.Context, email account.Email) (*account.Account, error) {
	var acc account.Account
	err := repo.db.Model(ctx, &acc).Where("email = ?", email).First(&acc).Error
	if err != nil {
		return nil, errors.Wrap(err, "find account by email")
	}

	return &acc, nil
}

// List all accounts.
func (repo *AccountRepository) List(ctx context.Context) ([]*account.Account, error) {
	var accounts []*account.Account
	err := repo.db.Model(ctx, &accounts).Find(&accounts).Error
	if err != nil {
		return nil, errors.Wrap(err, "find accounts")
	}

	return accounts, nil
}

// Create inserts a new Account into the repository.
func (repo *AccountRepository) Create(ctx context.Context, acc *account.Account) error {
	exists, err := repo.checkEmailExists(ctx, acc.Email)
	if err != nil {
		return errors.Wrap(err, "check email exists")
	}

	if exists {
		return account.ErrEmailAlreadyExists
	}

	err = repo.db.Model(ctx, acc).Create(acc).Error
	if err != nil {
		return errors.Wrap(err, "create account")
	}

	return nil
}

// Update updates an existing account in the repository.
func (repo *AccountRepository) Update(ctx context.Context, acc *account.Account) error {
	err := repo.db.Model(ctx, acc).Save(acc).Error
	if err != nil {
		return errors.Wrap(err, "update account")
	}

	return nil
}

func (repo *AccountRepository) checkEmailExists(ctx context.Context, email account.Email) (bool, error) {
	var count int
	err := repo.db.Model(ctx, &account.Account{}).
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		return false, errors.Wrap(err, "query accounts with email count")
	}

	return count > 0, nil
}
