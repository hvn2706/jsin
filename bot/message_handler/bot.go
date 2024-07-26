package message_handler

import (
	"context"
	"jsin/database"
	"jsin/external/s3"
	"jsin/logger"
)

type IMessageHandler interface {
	HandleMessage(ctx context.Context, message string) (MessageDTO, error)
	randomImageCmd(ctx context.Context) (MessageDTO, error)
}

type MessageHandler struct {
	s3client s3.IClient
	gdb      database.DBAdapter
}

func NewMessageHandler(s3client s3.IClient) IMessageHandler {
	return &MessageHandler{
		s3client: s3client,
		gdb:      database.GetDBInstance(),
	}
}

func (b *MessageHandler) HandleMessage(ctx context.Context, message string) (MessageDTO, error) {
	if message == jsinCommand {
		return b.randomImageCmd(ctx)
	}

	return MessageDTO{}, nil
}

func (b *MessageHandler) randomImageCmd(ctx context.Context) (MessageDTO, error) {
	var randImageKey string

	err := b.gdb.DB().Table("image").
		Select("file_name").
		Where("nsfw = ?", false).
		Order("rand()").
		Limit(1).
		Find(&randImageKey).Error
	if err != nil {
		return MessageDTO{}, err
	}

	img, err := b.s3client.GetImage(ctx, randImageKey)
	if err != nil {
		logger.Errorf("===== Get image failed: %+v", err.Error())
		return MessageDTO{}, err
	}
	return MessageDTO{
		Message: randImageKey,
		Object: &ObjectDTO{
			ObjectKey: randImageKey,
			Object:    img,
		},
	}, nil
}
