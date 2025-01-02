package message_handler

import (
	"context"
	"strings"

	"github.com/urfave/cli"

	"jsin/config"
	"jsin/external/Custom3rdParties"
	"jsin/external/s3"
	"jsin/logger"
	"jsin/pkg/common"
	"jsin/pkg/constants"
	"jsin/pkg/storage"
)

type IMessageHandler interface {
	HandleMessage(ctx context.Context, message string) (*MessageDTO, error)
	RandomImageCron(ctx context.Context) (*MessageDTO, error)
}

type MessageHandler struct {
	config                 config.Config
	s3client               s3.IClient
	custom3rdPartiesClient Custom3rdParties.IClient
	imageStorage           storage.ImageStorage
	cronJobStorage         storage.CronJobStorage
}

func NewMessageHandler(cfg config.Config) IMessageHandler {
	s3client := s3.NewClient(cfg.ExternalService.S3)
	custom3rdPartiesClient := Custom3rdParties.NewClient(cfg.ExternalService.Custom3rdParties)
	return &MessageHandler{
		config:                 cfg,
		s3client:               s3client,
		custom3rdPartiesClient: custom3rdPartiesClient,
		imageStorage:           storage.NewImageStorage(),
		cronJobStorage:         storage.NewCronJobStorage(),
	}
}

// HandleMessage generates content based on the message received
func (b *MessageHandler) HandleMessage(ctx context.Context, message string) (*MessageDTO, error) {
	args := strings.Split(message, " ")
	if args[0] != jsinCommand {
		return nil, nil
	}

	var generatedContent *MessageDTO
	var err error
	handler := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "type, t",
			},
		},
		Action: func(ctxCLI *cli.Context) error {
			generatedContent, err = b.randomImageCmd(ctx, ctxCLI.String("type"))
			return err
		},
		Commands: []cli.Command{
			{
				Name: "help",
				Action: func(ctxCLI *cli.Context) error {
					generatedContent, err = b.generateHelpContent(ctx)
					return err
				},
			},
			{
				Name: "cron",
				Action: func(ctxCLI *cli.Context) error {
					cronType := args[len(args)-2]
					cronTime := args[len(args)-1]
					err := common.IsValidTimeFormat(cronTime)
					if err != nil {
						return err
					}

					if cronType != constants.DailyType {
						// TODO: Add support for other cron types
						generatedContent = &MessageDTO{
							Message: "Currently only support daily",
						}
						return nil
					}

					generatedContent, err = b.generateCronJob(ctx, cronType, common.ConvertToCronFormat(cronTime))
					return err
				},
			},
			{
				Name: b.config.ExternalService.Custom3rdParties.Command,
				Action: func(ctxCLI *cli.Context) error {
					generatedContent, err = b.randomImageFrom3rdParties(ctx)
					return err
				},
			},
		},
	}

	if err := handler.Run(args); err != nil {
		logger.Errorf("===== Run command failed: %+v", err.Error())
		return nil, err
	}

	return generatedContent, nil
}

func (b *MessageHandler) RandomImageCron(ctx context.Context) (*MessageDTO, error) {
	randImageKey, err := b.imageStorage.RandomImage(ctx, "")
	if err != nil {
		logger.Errorf("===== Get random image failed: %+v", err.Error())
		return nil, err
	}

	img, err := b.s3client.GetImage(ctx, randImageKey)
	if err != nil {
		logger.Errorf("===== Get image failed: %+v", err.Error())
		return nil, err
	}

	return &MessageDTO{
		Message: randImageKey,
		Object: &ObjectDTO{
			ObjectKey: randImageKey,
			Object:    img,
		},
	}, nil
}

func (b *MessageHandler) randomImageCmd(ctx context.Context, imgType string) (*MessageDTO, error) {
	randImageKey, err := b.imageStorage.RandomImage(ctx, imgType)
	if err != nil {
		logger.Errorf("===== Get random image failed: %+v", err.Error())
		return nil, err
	}

	img, err := b.s3client.GetImage(ctx, randImageKey)
	if err != nil {
		logger.Errorf("===== Get image failed: %+v", err.Error())
		return nil, err
	}

	return &MessageDTO{
		Message: randImageKey,
		Object: &ObjectDTO{
			ObjectKey: randImageKey,
			Object:    img,
		},
	}, nil
}

func (b *MessageHandler) generateHelpContent(ctx context.Context) (*MessageDTO, error) {
	content := b.config.HelpContent

	return &MessageDTO{
		Message: content,
	}, nil
}

func (b *MessageHandler) generateCronJob(
	ctx context.Context,
	cronJobType string,
	cronTime string,
) (*MessageDTO, error) {
	content := b.config.TelegramBot.CreatCronJobContent

	_, err := b.cronJobStorage.AddCronJob(ctx, cronTime, cronJobType)
	if err != nil {
		return nil, err
	}
	return &MessageDTO{
		Message: content,
	}, nil
}

func (b *MessageHandler) randomImageFrom3rdParties(ctx context.Context) (*MessageDTO, error) {
	img, err := b.custom3rdPartiesClient.GetRandomImageFrom3rdParties(ctx)
	if err != nil {
		return nil, err
	}

	return &MessageDTO{
		Message: b.config.ExternalService.Custom3rdParties.Command,
		Object: &ObjectDTO{
			Object: img,
		},
	}, nil
}
