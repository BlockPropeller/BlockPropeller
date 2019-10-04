package auth

import (
	"chainup.dev/chainup"
	"chainup.dev/lib/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func whoamiCmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "whoami",
		Usage: "Return the currently authenticated account.",
		Action: func(c *cli.Context) {
			token, err := getToken()
			if errors.Cause(err) == errTokenNotFound {
				log.Info("not logged in")
				return
			}
			if err != nil {
				log.ErrorErr(err, "failed getting account token")
				return
			}

			acc, err := app.AccountService.Authenticate(token)
			if err != nil {
				log.ErrorErr(err, "failed authenticating user")
				return
			}

			log.Info("logged in", log.Fields{
				"id":    acc.ID,
				"email": acc.Email,
				"token": token,
			})
		},
	}
}
