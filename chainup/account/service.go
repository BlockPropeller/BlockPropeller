package account

import (
	"context"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidPassword is returned when we cannot match a provided password with the stored hash.
	ErrInvalidPassword = errors.New("invalid password")
)

// Service is responsible for creating and manipulating the Accounts.
type Service struct {
	accRepo  Repository
	tokenSvc *TokenService
}

// NewService returns a new Service instance.
func NewService(accRepo Repository, tokenSvc *TokenService) *Service {
	return &Service{accRepo: accRepo, tokenSvc: tokenSvc}
}

// Register an Account with the platform.
func (s *Service) Register(email Email, password ClearPassword) (*Account, Token, error) {
	if err := email.Validate(); err != nil {
		return nil, NilToken, errors.Wrap(err, "invalid email")
	}
	if err := password.Validate(); err != nil {
		return nil, NilToken, errors.Wrap(err, "invalid password")
	}

	pass, err := GeneratePassword(password)
	if err != nil {
		return nil, NilToken, errors.Wrap(err, "generate password")
	}

	acc := NewAccount(email, pass)

	err = s.accRepo.Create(context.TODO(), acc)
	if err != nil {
		return nil, NilToken, errors.Wrap(err, "create account")
	}

	token, err := s.tokenSvc.GenerateToken(acc.ID)
	if err != nil {
		return nil, NilToken, errors.Wrap(err, "generate token")
	}

	return acc, token, nil
}

// Login to an Account using a token.
func (s *Service) Login(email Email, password ClearPassword) (Token, error) {
	acc, err := s.accRepo.FindByEmail(context.TODO(), email)
	if err != nil {
		return NilToken, err
	}

	if err = acc.Password.Compare(password); err != nil {
		return NilToken, ErrInvalidPassword
	}

	token, err := s.tokenSvc.GenerateToken(acc.ID)
	if err != nil {
		return NilToken, err
	}

	return token, nil
}

// Authenticate an access token as an Account.
func (s *Service) Authenticate(token Token) (*Account, error) {
	id, err := s.tokenSvc.ParseToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse token")
	}

	acc, err := s.accRepo.FindByID(context.TODO(), id)
	if err != nil {
		return nil, errors.Wrap(err, "find account by id")
	}

	return acc, nil
}

// ChangePassword changes the password of an Account.
func (s *Service) ChangePassword(acc *Account, oldPassword ClearPassword, newPassword ClearPassword) error {
	if err := acc.Password.Compare(oldPassword); err != nil {
		return errors.Wrap(err, "invalid old password")
	}

	password, err := GeneratePassword(newPassword)
	if err != nil {
		return errors.Wrap(err, "generate new password")
	}

	acc.Password = password

	err = s.accRepo.Update(context.TODO(), acc)
	if err != nil {
		return errors.Wrap(err, "update account")
	}

	return nil
}
