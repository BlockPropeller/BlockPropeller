package auth

import (
	"blockpropeller.dev/blockpropeller/cmd/blockctl/util/localauth"
	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func logoutCmd() cli.Command {
	return cli.Command{
		Name:  "logout",
		Usage: "Logout from a BlockPropeller account.",
		Action: func(c *cli.Context) {
			err := localauth.DeleteToken()
			if errors.Cause(err) == localauth.ErrTokenNotFound {
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
