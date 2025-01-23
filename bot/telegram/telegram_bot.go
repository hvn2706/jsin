package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

	"jsin/bot/cron_handler"
	"jsin/bot/message_handler"
	"jsin/config"
	"jsin/logger"
	"jsin/pkg/constants"
)

type ITelegramBot interface {
	Serve() error
}

type Bot struct {
	cfg           config.TelegramBot
	botHandler    message_handler.IMessageHandler
	cronHandler   cron_handler.CronHandler
	bot           *tgbotapi.BotAPI
	cronScheduler *cron.Cron
}

func NewTelegramBot(cfg config.Config) ITelegramBot {
	botHandler := message_handler.NewMessageHandler(cfg)
	cronHandler := cron_handler.NewCronHandler()
	cronScheduler := cron.New()
	return &Bot{
		cfg:           cfg.TelegramBot,
		botHandler:    botHandler,
		cronHandler:   cronHandler,
		cronScheduler: cronScheduler,
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

	ctx := context.Background()
	go func() {
		err = b.sendImageCron(ctx)
		if err != nil {
			logger.Errorf("===== Send image cron failed: %+v", err.Error())
			return
		}
	}()

	u := tgbotapi.NewUpdate(b.cfg.Offset)
	u.Timeout = b.cfg.Timeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			ctx := context.WithValue(
				context.Background(),
				constants.ChatIDKey,
				fmt.Sprintf("%d", update.Message.Chat.ID),
			)
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
				err = b.sendImage(update, *generateContent.Object)
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
