package main

import (
	"github.com/urfave/cli"
	"jsin/config"
	"jsin/logger"
	"jsin/server"
	"log"
	"os"
)

func main() {
	// 1. Load config
	config.Load()

	// 2. Init logger
	_, err := logger.InitLogger(config.GlobalCfg.Logger)
	if err != nil {
		log.Fatalf("===== Init logger failed: %+v", err.Error())
	}

	// 3. Init server
	app := &cli.App{
		Name:  "jsin",
		Usage: "jsin provides you heaven",
		Action: func(ctx *cli.Context) error {
			srv := server.NewServer(config.GlobalCfg)
			return srv.Serve(config.GlobalCfg.Server.HTTP)
		},
	}

	// 4. Run server
	if err = app.Run(os.Args); err != nil {
		logger.Fatalf("===== Run server failed: %+v", err.Error())
	}
}
