package auth

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"chainup.dev/chainup/account"
	"github.com/pkg/errors"
)

var (
	errTokenNotFound = errors.New("token not found")
)

func getToken() (account.Token, error) {
	tokenFile, err := getTokenFile()
	if err != nil {
		return account.NilToken, errors.Wrap(err, "get token file")
	}

	_, err = os.Stat(tokenFile)
	if os.IsNotExist(err) {
		return account.NilToken, errTokenNotFound
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

func setToken(token account.Token) error {
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

func deleteToken() error {
	tokenFile, err := getTokenFile()
	if err != nil {
		return errors.Wrap(err, "get token file")
	}

	_, err = os.Stat(tokenFile)
	if os.IsNotExist(err) {
		return errTokenNotFound
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

	configDir := filepath.Join(homeDir, ".chainup")

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "create ChainUP config dir")
	}

	tokenFile := filepath.Join(configDir, "jwt_token")

	return tokenFile, nil
}
