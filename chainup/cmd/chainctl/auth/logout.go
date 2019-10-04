package auth

import (
	"chainup.dev/lib/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func logoutCmd() cli.Command {
	return cli.Command{
		Name:  "logout",
		Usage: "Logout from a ChainUP account.",
		Action: func(c *cli.Context) {
			err := deleteToken()
			if errors.Cause(err) == errTokenNotFound {
				log.Info("already logged out")
				return
			}
			if err != nil {
				log.ErrorErr(err, "failed deleting token")
				return
			}

			log.Info("successfully logged out")
		},
	}
}
