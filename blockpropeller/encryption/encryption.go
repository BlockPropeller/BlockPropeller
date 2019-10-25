package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
)

var secretKey string

// Init the encryption module.
func Init(s string) {
	secretKey = s
}

// Encrypt sensitive data and return the encrypted result.
func Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return nil, errors.Wrap(err, "new cypher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "new gcm cypher")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "read nonce")
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return []byte(encoded), nil
}

// Decrypt sensitive data from an encrypted result.
func Decrypt(data []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "base64 decode data")
	}

	block, err := aes.NewCipher(getKey())
	if err != nil {
		return nil, errors.Wrap(err, "new cypher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "new gcm cypher")
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get plaintext")
	}

	return plaintext, nil
}

func getKey() []byte {
	if secretKey == "" {
		panic("encryption module not initialized")
	}

	return []byte(createHash(secretKey))
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
