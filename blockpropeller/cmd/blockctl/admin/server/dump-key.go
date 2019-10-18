package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/lib/log"
	"github.com/urfave/cli"
)

func dumpKeyCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "dump-key",
		Usage: "Dump the private key used for accessing the server",
		Action: func(c *cli.Context) {
			if !c.Args().Present() {
				log.Error("please enter a server ID")
				return
			}

			ctx := context.Background()

			srvID := infrastructure.ServerIDFromString(c.Args().First())

			srv, err := app.ServerRepository.Find(ctx, srvID)
			if err != nil {
				log.ErrorErr(err, "failed listing servers")
				return
			}

			privKey := srv.SSHKey.EncodedPrivateKey()

			keyPath := filepath.Join(os.TempDir(), "keys", srv.ID.String()+".pkey")
			err = os.MkdirAll(filepath.Dir(keyPath), 0755)
			if err != nil {
				log.ErrorErr(err, "failed creating key dir")
				return
			}

			err = ioutil.WriteFile(keyPath, []byte(privKey), 0400)
			if err != nil {
				log.ErrorErr(err, "failed saving private key")
				return
			}

			log.Info("saved private key", log.Fields{
				"cmd": fmt.Sprintf("ssh root@%s -i %s", srv.IPAddress.String(), keyPath),
			})
		},
	}
}
