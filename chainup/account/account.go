package account

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	// NilID is the value of a nil account ID.
	NilID ID
)

// ID is a unique server identifier.
type ID string

// NewID returns a new unique ID.
func NewID() ID {
	return ID(uuid.NewV4().String())
}

// IDFromString converts the provided string to an ID.
func IDFromString(id string) ID {
	return ID(id)
}

// String satisfies the Stringer interface.
func (id ID) String() string {
	return string(id)
}

// Account represents an identity on a ChainUP platform.
type Account struct {
	ID ID `json:"id" gorm:"type:varchar(36) not null"`

	Email    Email    `json:"email" gorm:"type:varchar(255) not null;unique_index"`
	Password Password `json:"-" gorm:"type:varchar(255)"`

	CreatedAt time.Time `json:"created_at" gorm:"type:datetime not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime not null;default:CURRENT_TIMESTAMP"`
}

// NewAccount returns a new Account instance.
func NewAccount(email Email, password Password) *Account {
	return &Account{
		ID:       NewID(),
		Email:    email,
		Password: password,
	}
}
