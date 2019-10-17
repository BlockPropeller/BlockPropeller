package auth

import (
	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/cmd/blockctl/util/localauth"
	"blockpropeller.dev/lib/log"
	"github.com/urfave/cli"
)

func loginCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "Login with a BlockPropeller Account",
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

			token, err := app.AccountService.Login(email, password)
			if err != nil {
				log.ErrorErr(err, "register new account")
				return
			}

			err = localauth.SetToken(token)
			if err != nil {
				log.ErrorErr(err, "failed saving token")
				return
			}

			log.Info("successfully logged in", log.Fields{
				"token": token,
			})
		},
	}
}
