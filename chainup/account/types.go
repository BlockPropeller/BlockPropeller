package account

import (
	"strings"

	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// ClearPassword is a password in a plain text format.
type ClearPassword string

// NewClearPassword returns a new ClearPassword instance.
func NewClearPassword(password string) ClearPassword {
	return ClearPassword(password)
}

// Validate checks that the ClearPassword matches the minimum requirements for a password.
func (cp ClearPassword) Validate() error {
	if len(cp) < 6 {
		return errors.New("password must have at least 6 characters")
	}

	return nil
}

// String satisfies the stringer interface.
func (cp ClearPassword) String() string {
	return string(cp)
}

// Password is the hash of an account password.
type Password string

// GeneratePassword returns a Password given a raw password string.
func GeneratePassword(password ClearPassword) (Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", errors.Wrap(err, "generate password hash")
	}

	return Password(hash), nil
}

// Compare check if the given password matches the password hash.
func (p Password) Compare(password ClearPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
}

// Email is an Accounts email.
type Email string

// NewEmail returns a new Email instance.
func NewEmail(email string) Email {
	return Email(strings.ToLower(email))
}

// Validate checks that an email is valid.
func (e Email) Validate() error {
	err := checkmail.ValidateFormat(e.String())
	if err != nil {
		return errors.Wrap(err, "validate email")
	}

	return nil
}

// String satisfies the Stringer interface.
func (e Email) String() string {
	return string(e)
}
