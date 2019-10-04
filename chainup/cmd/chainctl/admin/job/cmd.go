package job

import (
	"chainup.dev/chainup"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for running and operating jobs.
func Cmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "job",
		Usage: "Job related commands",
		Subcommands: []cli.Command{
			listCmd(app),
			runCmd(app),
		},
	}
}
