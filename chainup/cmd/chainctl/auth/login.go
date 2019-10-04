package auth

import (
	"chainup.dev/chainup"
	"chainup.dev/chainup/account"
	"chainup.dev/lib/log"
	"github.com/urfave/cli"
)

func loginCmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "Login with a ChainUP Account",
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

			err = setToken(token)
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
