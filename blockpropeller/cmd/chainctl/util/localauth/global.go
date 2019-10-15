package localauth

import (
	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/lib/log"
)

// Account is an authenticated account based on locally stored token.
var Account *account.Account

// Authenticate checks if there is an authenticated token and if so,
// sets the authenticated user to be globally accessible for future callers.
func Authenticate(app *blockpropeller.App) {
	token, err := GetToken()
	if err != nil {
		log.Debug("failed getting auth token from local storage", log.Fields{
			"err": err,
		})
		return
	}

	acc, err := app.AccountService.Authenticate(token)
	if err != nil {
		log.Debug("failed authenticating account from token", log.Fields{
			"token": token.String(),
			"err":   err,
		})
		return
	}

	Account = acc
}
