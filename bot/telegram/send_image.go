package telegram

import (
	"context"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

	"jsin/bot/message_handler"
	"jsin/logger"
	"jsin/pkg/constants"
)

func (b *Bot) sendImageCron(ctx context.Context) error {
	for {
		b.cronScheduler.Stop()

		b.cronScheduler = cron.New()
		cronJobs, err := b.cronHandler.ListCronJobs(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch cron jobs: %v", err)
			return err
		}

		for _, job := range cronJobs {
			chatID, err := strconv.ParseInt(job.ChatID, 10, 64)
			if err != nil {
				logger.Errorf("Invalid chat ID: %v", err)
				continue
			}

			_, err = b.cronScheduler.AddFunc(job.CronJob, func() {
				err := b.sendImageByCron(chatID)
				if err != nil {
					return
				}
			})
			if err != nil {
				logger.Errorf("Error scheduling cron job for chat ID %d: %v", chatID, err)
			}
		}

		b.cronScheduler.Start()

		select {
		case <-ctx.Done():
			b.cronScheduler.Stop()
			logger.Info("Scheduler stopped")
			return err
		case <-time.After(constants.IntervalTime * time.Second):
		}
	}
}

func (b *Bot) sendImage(update tgbotapi.Update, object message_handler.ObjectDTO) error {
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

func (b *Bot) sendImageByObject(chatID int64, object message_handler.ObjectDTO) error {
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

func (b *Bot) sendImageByCron(chatID int64) error {
	currentTime := time.Now().Format(constants.DayFormater)
	generateContent, err := b.botHandler.RandomImageCron(context.Background())
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, constants.Sorry)
		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	}
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  generateContent.Object.ObjectKey,
		Bytes: generateContent.Object.Object,
	})
	photo.Caption = b.cfg.CronJobImageCaption

	_, err = b.bot.Send(photo)
	if err != nil {
		logger.Errorf("===== Send image failed for message ID %v: %+v ", chatID, err.Error())
		return err
	}

	logger.Infof("Your gift was send successfully at %s", currentTime)
	return nil
}
