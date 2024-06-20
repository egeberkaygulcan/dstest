package network

// Message is a message that can be sent between processes
type Message struct {
	// Sender is the id of the sender process.
	Sender int
	// Receiver is the id of the receiver process.
	Receiver int
	// Payload is the data that is being sent.
	Payload any
	// Message protocol type
	Type MessageType
	// Name of the action message carries
	Name string
	// Unique message Id
	MessageId uint64
	// Metadata is any additional data that can be annotated to the message.
	//Metadata map[string]any

	// A channel to trigger sending the response
	Send chan struct{}
}

type MessageType string

const (
	GRPC MessageType = "GRPC"		
)


