package message_handler

import (
	"context"

	"jsin/external/s3"
	"jsin/logger"
	"jsin/pkg/storage"
)

type IMessageHandler interface {
	HandleMessage(ctx context.Context, message string) (*MessageDTO, error)

	randomImageCmd(ctx context.Context) (*MessageDTO, error)
}

type MessageHandler struct {
	s3client     s3.IClient
	imageStorage storage.ImageStorage
}

func NewMessageHandler(s3client s3.IClient) IMessageHandler {
	return &MessageHandler{
		s3client:     s3client,
		imageStorage: storage.NewImageStorage(),
	}
}

// HandleMessage generates content based on the message received
func (b *MessageHandler) HandleMessage(ctx context.Context, message string) (*MessageDTO, error) {
	if message == jsinCommand {
		return b.randomImageCmd(ctx)
	}

	return nil, nil
}

func (b *MessageHandler) randomImageCmd(ctx context.Context) (*MessageDTO, error) {
	randImageKey, err := b.imageStorage.RandomImage(ctx, false)
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
