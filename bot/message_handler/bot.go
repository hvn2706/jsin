package message_handler

type IMessageHandler interface {
	HandleMessage(message string) (string, error)
}

type MessageHandler struct {
}

func NewMessageHandler() IMessageHandler {
	return &MessageHandler{}
}

func (b *MessageHandler) HandleMessage(message string) (string, error) {
	return message, nil
}
