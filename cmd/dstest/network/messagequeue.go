package network

import (
	"container/list"
)

// MessageQueue is a queue of messages
type MessageQueue struct {
	List list.List
}

// Message is a message that can be sent between processes
type Message struct {
	Sender   int
	Receiver int
	Payload  any
	Metadata map[string]any
}

// Init initializes the message queue
func (mq *MessageQueue) Init() {
	mq.List.Init()
}

// PushBack adds a message to the back of the queue
func (mq *MessageQueue) PushBack(m Message) {
	mq.List.PushBack(m)
}

// PopFront removes a message from the front of the queue
func (mq *MessageQueue) PopFront() Message {
	e := mq.List.Front()
	mq.List.Remove(e)
	return e.Value.(Message)
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
