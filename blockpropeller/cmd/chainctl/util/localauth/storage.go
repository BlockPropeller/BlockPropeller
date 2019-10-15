package localauth

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"blockpropeller.dev/blockpropeller/account"
	"github.com/pkg/errors"
)

var (
	// ErrTokenNotFound is an error returned when a token is not found locally.
	ErrTokenNotFound = errors.New("token not found")
)

// GetToken returns an account token stored locally if it exists.
func GetToken() (account.Token, error) {
	tokenFile, err := getTokenFile()
	if err != nil {
		return account.NilToken, errors.Wrap(err, "get token file")
	}

	_, err = os.Stat(tokenFile)
	if os.IsNotExist(err) {
		return account.NilToken, ErrTokenNotFound
	}
	if err != nil {
		return account.NilToken, errors.Wrap(err, "stat token file")
	}

	data, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return account.NilToken, errors.Wrap(err, "read token file")
	}

	return account.NewToken(string(data)), nil
}

// SetToken sets the provided Token to local storage for future authentication.
func SetToken(token account.Token) error {
	tokenFile, err := getTokenFile()
	if err != nil {
		return errors.Wrap(err, "get token file")
	}

	err = ioutil.WriteFile(tokenFile, []byte(token.String()), 0644)
	if err != nil {
		return errors.Wrap(err, "write JWT token to config dir")
	}

	return nil
}

// DeleteToken removes the Token from local storage if it exists.
func DeleteToken() error {
	tokenFile, err := getTokenFile()
	if err != nil {
		return errors.Wrap(err, "get token file")
	}

	_, err = os.Stat(tokenFile)
	if os.IsNotExist(err) {
		return ErrTokenNotFound
	}
	if err != nil {
		return errors.Wrap(err, "stat token file")
	}

	err = os.Remove(tokenFile)
	if err != nil {
		return errors.Wrap(err, "delete token file")
	}

	return nil
}

func getTokenFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "get user home dir")
	}

	configDir := filepath.Join(homeDir, ".blockpropeller")

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "create BlockPropeller config dir")
	}

	tokenFile := filepath.Join(configDir, "jwt_token")

	return tokenFile, nil
}
