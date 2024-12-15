package message_handler

import (
	"context"
	"strings"

	"github.com/urfave/cli"

	"jsin/config"
	"jsin/external/s3"
	"jsin/logger"
	"jsin/pkg/storage"
)

type IMessageHandler interface {
	HandleMessage(ctx context.Context, message string) (*MessageDTO, error)

	randomImageCmd(ctx context.Context, imgType string) (*MessageDTO, error)
	generateHelpContent(ctx context.Context) (*MessageDTO, error)
}

type MessageHandler struct {
	config       config.Config
	s3client     s3.IClient
	imageStorage storage.ImageStorage
}

func NewMessageHandler(cfg config.Config) IMessageHandler {
	s3client := s3.NewClient(cfg.ExternalService.S3)
	return &MessageHandler{
		config:       cfg,
		s3client:     s3client,
		imageStorage: storage.NewImageStorage(),
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
		},
	}

	if err := handler.Run(args); err != nil {
		logger.Errorf("===== Run command failed: %+v", err.Error())
		return nil, err
	}

	return generatedContent, nil
}

func (b *MessageHandler) randomImageCmd(ctx context.Context, imgType string) (*MessageDTO, error) {
	logger.Infof("img type: %s", imgType)
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
