package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"jsin/bot/message_handler"
	"jsin/config"
	"jsin/external/s3"
	"jsin/logger"
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
	botHandler := message_handler.NewMessageHandler(s3.NewClient(cfg.ExternalService.S3))
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

	u := tgbotapi.NewUpdate(b.cfg.Offset)
	u.Timeout = b.cfg.Timeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			generateContent, err := b.botHandler.HandleMessage(context.Background(), update.Message.Text)
			if err != nil {
				logger.Errorf("===== Handle message failed: %+v", err.Error())
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hmm, something went wrong"))
				continue
			}
			if generateContent.Message == "" {
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

func (b *Bot) SendImage(update tgbotapi.Update, object message_handler.ObjectDTO) error {
	file := tgbotapi.FileBytes{
		Name:  object.ObjectKey,
		Bytes: object.Object,
	}
	message, err := b.bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, file))
	if err != nil {
		logger.Errorf("===== Send image failed: %+v", err.Error())
		return err
	}
	logger.Infof("Image sent, message id: %d", message.MessageID)

	return nil
}
