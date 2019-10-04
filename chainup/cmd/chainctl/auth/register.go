package auth

import (
	"chainup.dev/chainup"
	"chainup.dev/chainup/account"
	"chainup.dev/lib/log"
	"github.com/urfave/cli"
)

func registerCmd(app *chainup.App) cli.Command {
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

			acc, err := app.AccountService.Register(email, password)
			if err != nil {
				log.ErrorErr(err, "register new account")
				return
			}

			log.Info("created new account", log.Fields{
				"id":    acc.ID,
				"email": acc.Email,
			})
		},
	}
}
