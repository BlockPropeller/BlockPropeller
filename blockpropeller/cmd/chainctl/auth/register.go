package auth

import (
	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/util/localauth"
	"blockpropeller.dev/lib/log"
	"github.com/urfave/cli"
)

func registerCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "register",
		Usage: "Register a new account",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "email",
				Usage:    "Account email",
				Required: true,
			},
			cli.StringFlag{
				Name:     "password",
				Usage:    "Account password",
				Required: true,
			},
		},
		Action: func(c *cli.Context) {
			email := account.NewEmail(c.String("email"))
			password := account.NewClearPassword(c.String("password"))

			acc, token, err := app.AccountService.Register(email, password)
			if err != nil {
				log.ErrorErr(err, "register new account")
				return
			}

			err = localauth.SetToken(token)
			if err != nil {
				log.ErrorErr(err, "failed saving token")
				return
			}

			log.Info("created new account", log.Fields{
				"id":    acc.ID,
				"email": acc.Email,
			})
		},
	}
}
