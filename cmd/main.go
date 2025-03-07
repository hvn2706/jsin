package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/user"
	"github.com/urfave/cli"

	"jsin/bot/telegram"
	"jsin/cmd/job"
	"jsin/config"
	"jsin/database"
	"jsin/logger"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	client, err := bucketeer.NewSDK(
		ctx,
		bucketeer.WithAPIKey("test"),
		bucketeer.WithHost("test"),
		bucketeer.WithTag("test"),
		bucketeer.WithEnableLocalEvaluation(true),          // <--- Enable the local evaluation
		bucketeer.WithCachePollingInterval(10*time.Minute), // <--- Change the default interval if needed
	)
	if err != nil {
		log.Fatalf("Failed initialize the new client: %v", err)
	}

	jsinUser := user.NewUser(
		"END_USER_ID",
		nil, // The jsinUser attributes are optional
	)
	showNewFeature := client.BoolVariation(ctx, jsinUser, "feature-go-server-e2e-string", false)
	if showNewFeature {
		// The Application code to show the new feature
	} else {
		// The code to run when the feature is off
	}

	// 4. Init server
	app := &cli.App{
		Commands: []cli.Command{
			{
				Name:  "jsin-telegram",
				Usage: "jsin provides you heaven telegram bot",
				Action: func(ctx *cli.Context) error {
					bot := telegram.NewTelegramBot(config.GlobalCfg)
					err = bot.Serve()
					if err != nil {
						logger.Errorf("===== Run telegram bot failed: %+v", err.Error())
						return err
					}
					return nil
				},
			},
			{
				Name:  "jsin-migration",
				Usage: "migrate object to s3 and save url to db",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "type",
					},
				},
				Action: func(ctx *cli.Context) error {
					err = job.StartMigrationObjectJob(context.Background(), ctx.String("type"))
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	// 4. Run server
	if err = app.Run(os.Args); err != nil {
		logger.Fatalf("===== Run server failed: %+v", err.Error())
	}
}
