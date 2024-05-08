package cli

import (
	"os"

	"github.com/urfave/cli/v2"

	li "github.com/predixus/analytics_framework/internal/logger"
)

type CliArguments struct {
	InitialiseDB bool
	Continue     bool
}

func ParseInputs() CliArguments {
	var cliArgs CliArguments

	app := &cli.App{
		Name:  "calc",
		Usage: "To create scalable and robust analytics pipelines",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "initDB",
				Value: false,
				Usage: "For provisioning a local Postgres DB",
			},
			&cli.BoolFlag{
				Name:  "continue",
				Value: false,
				Usage: "To continue to launching the platform after performing setup tasks",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cliArgs = CliArguments{
				InitialiseDB: cCtx.Bool("initDB"),
				Continue:     cCtx.Bool("continue"),
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		li.Logger.Fatal(err)
	}

	return cliArgs
}
