package message_handler

const (
	jsinCommand = "j.sin"
)

type MessageDTO struct {
	Message string
	Object  *ObjectDTO
}

type ObjectDTO struct {
	ObjectKey string
	Object    []byte
}
