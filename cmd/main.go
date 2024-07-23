package main

import (
	"github.com/urfave/cli"
	"jsin/bot/telegram"
	"jsin/config"
	"jsin/database"
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

	// 3. Init db
	err = database.GetDBInstance().Open(config.GlobalCfg.Database.MySQLConfig)
	if err != nil {
		logger.Fatalf("===== Open db failed: %+v", err.Error())
	}

	// 4. Init server
	app := &cli.App{
		Commands: []cli.Command{
			{
				Name:  "jsin-api",
				Usage: "jsin provides you heaven api",
				Action: func(ctx *cli.Context) error {
					srv := server.NewServer(config.GlobalCfg)
					return srv.Serve(config.GlobalCfg.Server.HTTP)
				},
			},
			{
				Name:  "jsin-telegram",
				Usage: "jsin provides you heaven telegram bot",
				Action: func(ctx *cli.Context) {
					bot := telegram.NewTelegramBot(config.GlobalCfg.TelegramBot)
					err = bot.Serve()
					if err != nil {
						logger.Fatalf("===== Run telegram bot failed: %+v", err.Error())
					}
				},
			},
		},
	}

	// 4. Run server
	if err = app.Run(os.Args); err != nil {
		logger.Fatalf("===== Run server failed: %+v", err.Error())
	}
}
