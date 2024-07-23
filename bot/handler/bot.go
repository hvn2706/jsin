package handler

type IBotHandler interface {
	HandleMessage(message string) (string, error)
}

type BotHandler struct {
}

func NewBotHandler() IBotHandler {
	return &BotHandler{}
}

func (b *BotHandler) HandleMessage(message string) (string, error) {
	return message, nil
}
