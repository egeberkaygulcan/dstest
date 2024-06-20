package network

import (
	"container/list"
	"log"
	"sync"
)

// MessageQueue is a queue of messages
type MessageQueue struct {
	mu   sync.Mutex
	List list.List
}

func (m *Message) SendMessage() {
	close(m.Send)
}

// Init initializes the message queue
func (mq *MessageQueue) Init() {
	mq.List.Init()
}

// PushBack adds a message to the back of the queue
func (mq *MessageQueue) PushBack(m *Message) {
	mq.mu.Lock()
	mq.List.PushBack(m)
	mq.mu.Unlock()
}

// PopFront removes a message from the front of the queue
func (mq *MessageQueue) PopFront() *Message {
	mq.mu.Lock()
	e := mq.List.Front()
	mq.List.Remove(e)
	value := e.Value.(*Message)
	mq.mu.Unlock()
	return value
}

// Len returns the length of the queue
func (mq *MessageQueue) Len() int {
	return mq.List.Len()
}

// Remove a specific message from the queue
func (mq *MessageQueue) Remove(m *Message) {
	mq.mu.Lock()
	// find element with value m
	for e := mq.List.Front(); e != nil; e = e.Next() {
		if e.Value == m {
			mq.List.Remove(e)
			break
		}
	}
	mq.mu.Unlock()
}

// Print
func (mq *MessageQueue) Print(Logger *log.Logger) {
	mq.mu.Lock()
	i := 0
	for e := mq.List.Front(); e != nil; e = e.Next() {
		m := e.Value.(Message)
		Logger.Printf("- [%d] Sender: %d, Receiver: %d, Payload: %s\n", i, m.Sender, m.Receiver, m.Payload)
		i += 1
	}
	mq.mu.Unlock()
}

func (mq *MessageQueue) Peek() *Message {
	var msg *Message = nil
	mq.mu.Lock()
	message := (mq.List.Front())
	if message != nil {
		msg = message.Value.(*Message)
	} else {
		msg = nil
	}
	mq.mu.Unlock()
	return msg
}
