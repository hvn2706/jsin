package telegram

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"jsin/bot/message_handler"
	"jsin/logger"
	"jsin/pkg/constants"
)

func (b *Bot) SendImage(update tgbotapi.Update, object message_handler.ObjectDTO) error {
	file := tgbotapi.FileBytes{
		Name:  object.ObjectKey,
		Bytes: object.Object,
	}
	message, err := b.bot.Send(
		tgbotapi.PhotoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           update.Message.Chat.ID,
					ReplyToMessageID: 0,
				},
				File: file,
			},
		},
	)
	if err != nil {
		logger.Errorf("===== Send image failed: %+v", err.Error())
		return err
	}
	logger.Infof("Image sent, message id: %d", message.MessageID)

	return nil
}

func (b *Bot) SendImageByObject(chatID int64, object message_handler.ObjectDTO) error {
	file := tgbotapi.FileBytes{
		Name:  object.ObjectKey,
		Bytes: object.Object,
	}

	message, err := b.bot.Send(
		tgbotapi.PhotoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID: chatID,
				},
				File: file,
			},
			Caption: constants.Greeting,
		},
	)

	if err != nil {
		logger.Errorf("===== Send image failed: %+v", err.Error())
		return err
	}

	logger.Infof("Image sent, message id: %d", message.MessageID)
	return nil
}

func (b *Bot) SendImageRandomDaily(chatID int) error {
	currentTime := time.Now().Format(constants.DayFormater)
	generateContent, err := b.botHandler.RandomImageCron(context.Background())

	if err != nil {
		msg := tgbotapi.NewMessage(int64(chatID), constants.Sorry)
		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	}
	photo := tgbotapi.NewPhoto(int64(chatID), tgbotapi.FileBytes{
		Name:  generateContent.Object.ObjectKey,
		Bytes: generateContent.Object.Object,
	})

	_, err = b.bot.Send(photo)
	if err != nil {
		logger.Errorf("===== Send image failed for message ID %s: %+v ", chatID, err.Error())
		return err
	}

	logger.Infof("Your gift was send successfully at %s", currentTime)
	return nil
}
