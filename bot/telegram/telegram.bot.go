package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"jsin/bot/handler"
	"jsin/config"
	"jsin/logger"
)

type ITelegramBot interface {
	Serve() error
}

type Bot struct {
	cfg        config.TelegramBot
	botHandler handler.IBotHandler
}

func NewTelegramBot(cfg config.TelegramBot) ITelegramBot {
	botHandler := handler.NewBotHandler()
	return &Bot{
		cfg:        cfg,
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
	logger.Infof("Telegram bot start to serve, bot name: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(b.cfg.Offset)
	u.Timeout = b.cfg.Timeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
