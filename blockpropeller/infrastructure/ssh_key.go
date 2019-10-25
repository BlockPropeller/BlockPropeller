package infrastructure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql/driver"
	"encoding/pem"

	"blockpropeller.dev/blockpropeller/encryption"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

const defaultKeySize = 4096

// PrivateKey is a database friendly representation of a rsa.PrivateKey.
type PrivateKey struct {
	rsa.PrivateKey
}

// NewPrivateKey returns a new PrivateKey instance.
func NewPrivateKey(pk *rsa.PrivateKey) *PrivateKey {
	return &PrivateKey{
		*pk,
	}
}

// Value implements the sql.Valuer interface.
func (pk *PrivateKey) Value() (driver.Value, error) {
	key := x509.MarshalPKCS1PrivateKey(&pk.PrivateKey)

	encrypted, err := encryption.Encrypt(key)
	if err != nil {
		return nil, errors.Wrap(err, "encrypt private key")
	}

	return encrypted, nil
}

// Scan implements the sql.Scanner interface.
func (pk *PrivateKey) Scan(src interface{}) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return errors.New("invalid private key type")
	}

	decrypted, err := encryption.Decrypt(data)
	if err != nil {
		return errors.Wrap(err, "decrypt private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(decrypted)
	if err != nil {
		return errors.Wrap(err, "parse private key")
	}

	*pk = *NewPrivateKey(privKey)

	return nil
}

// SSHKey is a private key that can be used as an authentication mechanism
// for logging into provisioned servers.
type SSHKey struct {
	Name       string      `json:"name" gorm:"type:varchar(100)"`
	PrivateKey *PrivateKey `json:"-" gorm:"type:text"`
}

// GenerateNewSSHKey generates a new random private key.
func GenerateNewSSHKey(name string) (*SSHKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultKeySize)
	if err != nil {
		return nil, errors.Wrap(err, "generate private key")
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate private key")
	}

	return &SSHKey{
		Name:       name,
		PrivateKey: NewPrivateKey(privateKey),
	}, nil
}

// EncodedPrivateKey encodes the private key into PEM format suitable for usage inside files.
func (key *SSHKey) EncodedPrivateKey() string {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(&key.PrivateKey.PrivateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return string(privatePEM)
}

// EncodedPublicKey encodes the public key into a format suitable for inclusion into authorized_keys file.
func (key *SSHKey) EncodedPublicKey() string {
	publicKey, err := ssh.NewPublicKey(&key.PrivateKey.PublicKey)
	if err != nil {
		panic(errors.Wrap(err, "create public key"))
	}

	return string(ssh.MarshalAuthorizedKey(publicKey))
}
