package network

import (
	"container/list"
	"log"
)

// MessageQueue is a queue of messages
type MessageQueue struct {
	List list.List
}

// Message is a message that can be sent between processes
type Message struct {
	// Sender is the id of the sender process.
	Sender int
	// Receiver is the id of the receiver process.
	Receiver int
	// Payload is the data that is being sent.
	Payload any
	// Metadata is any additional data that can be annotated to the message.
	//Metadata map[string]any

	// A channel to trigger sending the response
	Send chan struct{}
}

func (m *Message) send() {
	close(m.Send)
}

// Init initializes the message queue
func (mq *MessageQueue) Init() {
	mq.List.Init()
}

// PushBack adds a message to the back of the queue
func (mq *MessageQueue) PushBack(m *Message) {
	mq.List.PushBack(m)
}

// PopFront removes a message from the front of the queue
func (mq *MessageQueue) PopFront() *Message {
	e := mq.List.Front()
	mq.List.Remove(e)
	return e.Value.(*Message)
}

// Len returns the length of the queue
func (mq *MessageQueue) Len() int {
	return mq.List.Len()
}

// Remove a specific message from the queue
func (mq *MessageQueue) Remove(m *Message) {
	// find element with value m
	for e := mq.List.Front(); e != nil; e = e.Next() {
		if e.Value == m {
			mq.List.Remove(e)
			return
		}
	}
}

// Print
func (mq *MessageQueue) Print(Logger *log.Logger) {
	i := 0
	for e := mq.List.Front(); e != nil; e = e.Next() {
		m := e.Value.(Message)
		Logger.Printf("- [%d] Sender: %d, Receiver: %d, Payload: %s\n", i, m.Sender, m.Receiver, m.Payload)
		i += 1
	}
}
