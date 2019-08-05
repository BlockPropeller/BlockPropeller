package infrastructure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

const defaultKeySize = 4096

// SSHKey is a private key that can be used as an authentication mechanism
// for logging into provisioned servers.
type SSHKey struct {
	Name       string
	PrivateKey *rsa.PrivateKey
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
		PrivateKey: privateKey,
	}, nil
}

// EncodedPrivateKey encodes the private key into PEM format suitable for usage inside files.
func (key *SSHKey) EncodedPrivateKey() string {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(key.PrivateKey)

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
