package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/predixus/pdb_framework/internal/api"
	dlyr "github.com/predixus/pdb_framework/internal/datalayer"
	"github.com/predixus/pdb_framework/internal/grpc"
	li "github.com/predixus/pdb_framework/internal/logger"
	prov "github.com/predixus/pdb_framework/internal/provision_pg"
	setup "github.com/predixus/pdb_framework/internal/setup"
)

var InitDB, Continue bool

func parseInputs() {
	app := &cli.App{
		Name:  "pdb",
		Usage: "Initialise and run the Predixus DB (PDB)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "init-db",
				Value:       false,
				Usage:       "For provisioning a local Postgres DB",
				Destination: &InitDB,
			},
			&cli.BoolFlag{
				Name:        "continue",
				Value:       false,
				Usage:       "To continue to launching the platform after performing setup tasks",
				Destination: &Continue,
			},
		},
		Action: mainAction,
	}
	if err := app.Run(os.Args); err != nil {
		li.Logger.Fatal(err)
	}
}

func mainAction(ctx *cli.Context) error {
	if InitDB {
		li.Logger.Info("Initialising local DB")
		err := prov.Provision()
		if err != nil {
			return err
		}
		li.Logger.Info("Finished initialisation")
		if !Continue {
			return nil
		}
	}

	li.Logger.Info("Continuing to run the framework")
	db := &dlyr.Db{}
	api_server := &api.HttpServer{}
	grpc_server := &grpc.GrpcServer{}
	err := setup.Setup(db, grpc_server, api_server)
	return err
}

func main() {
	parseInputs()
}
