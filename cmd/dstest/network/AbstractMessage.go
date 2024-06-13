package network

type TranslatedMessage struct {
	Sender int
	Receiver int
	Name string
	// TODO - Additional metadata
}

type MessageType string

const (
	GRPC MessageType = "grpc"		
)

type AbstractMessage struct {
	Message *Message
	TranslatedMessage *TranslatedMessage
	Type MessageType // TODO - Create enum types for protocols
	MessageId uint64
}

