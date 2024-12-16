package telegram

import (
	"context"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

	"jsin/bot/message_handler"
	"jsin/config"
	"jsin/logger"
	"jsin/pkg/common"
	"jsin/pkg/constants"
	"jsin/pkg/storage"
)

type ITelegramBot interface {
	Serve() error
}

type Bot struct {
	cfg        config.TelegramBot
	botHandler message_handler.IMessageHandler
	bot        *tgbotapi.BotAPI
}

func NewTelegramBot(cfg config.Config) ITelegramBot {
	botHandler := message_handler.NewMessageHandler(cfg)
	return &Bot{
		cfg:        cfg.TelegramBot,
		botHandler: botHandler,
	}
}

func (b *Bot) Serve() error {
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		logger.Errorf("===== Init telegram bot failed: %+v", err.Error())
		return err
	}
	bot.Debug = b.cfg.Debug
	b.bot = bot
	logger.Infof("Telegram bot start to serve, bot name: %s", bot.Self.UserName)
	go b.Schedule(context.Background())

	u := tgbotapi.NewUpdate(b.cfg.Offset)
	u.Timeout = b.cfg.Timeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			logger.Infof("[%s-%d] %s",
				update.Message.From.UserName,
				update.Message.Chat.ID,
				update.Message.Text,
			)

			ctx := context.WithValue(context.Background(), "chatID", update.Message.Chat.ID)
			generateContent, err := b.botHandler.HandleMessage(ctx, update.Message.Text)
			if err != nil {
				logger.Errorf("===== Handle message failed: %+v", err.Error())
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				continue
			}
			if generateContent == nil || generateContent.Message == "" {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, generateContent.Message)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = bot.Send(msg)
			if err != nil {
				logger.Errorf("===== Send message failed: %+v", err.Error())
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hmm, something went wrong"))
				continue
			}

			if generateContent.Object != nil {
				err = b.SendImage(update, *generateContent.Object)
				if err != nil {
					logger.Errorf("===== Send image failed: %+v", err.Error())
					_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hmm, something went wrong"))
					continue
				}
			}
		}
	}

	return nil
}

func (b *Bot) Schedule(ctx context.Context) {
	c := cron.New(cron.WithLocation(common.LoadTimeZone()))

	for {
		c.Stop()
		c = cron.New(cron.WithLocation(common.LoadTimeZone()))

		cronStorage := storage.NewCronJobStorage()
		cronJobs, err := cronStorage.ListCronJobDaily(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch daily cron jobs: %v", err)
			continue
		}

		for _, job := range cronJobs {
			chatID, err := strconv.ParseInt(job.ChatID, 10, 64)
			if err != nil {
				logger.Errorf("Invalid chat ID: %v", err)
				continue
			}

			messageID, err := strconv.Atoi(job.ChatID)
			if err != nil {
				logger.Errorf("Invalid message ID: %v", err)
				continue
			}

			_, err = c.AddFunc(job.CronJob, func() {
				err := b.SendImageRandomDaily(messageID)
				if err != nil {
					return
				}
			})

			if err != nil {
				logger.Errorf("Error scheduling cron job for chat ID %d: %v", chatID, err)
			}
		}

		c.Start()
		logger.Info("Daily cron jobs scheduled successfully")

		select {
		case <-ctx.Done():
			c.Stop()
			logger.Info("Scheduler stopped")
			return
		case <-time.After(constants.IntervalTime * time.Second):
			logger.Info("Refreshing cron jobs...")
		}
	}
}
