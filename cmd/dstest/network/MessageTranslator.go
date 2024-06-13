package network

import (
	// "io"
	"log"
	"net/http/httputil"
	"os"
	"strings"
	"sync/atomic"
)

type MessageTranslator interface {
	Translate(Message) AbstractMessage
}

func NewMessageTranslator(messageType MessageType) MessageTranslator {
	switch messageType {
	case GRPC:
		return newGRPCTranslator()
	default:
		return nil
	}
}

type GRPCTranslator struct {
	index atomic.Uint64
	Log *log.Logger
}

func newGRPCTranslator() *GRPCTranslator {
	translator := new(GRPCTranslator)
	translator.Log = log.New(os.Stdout, "[GRPCTranslator]", log.LstdFlags)
	return translator
}

func (t *GRPCTranslator) Translate(message Message) AbstractMessage {
	payload := message.Payload.(Http2CPayload)

	abstractMessage := new(AbstractMessage)
	abstractMessage.Type = GRPC
	abstractMessage.MessageId = t.generateMessageId()
	abstractMessage.TranslatedMessage = new(TranslatedMessage)
	abstractMessage.TranslatedMessage.Sender = message.Sender
	abstractMessage.TranslatedMessage.Receiver = message.Receiver

	// buf := new(strings.Builder)
	// _, err := io.Copy(buf, payload.Request.Body)
	// if err != nil {
	// 	t.Log.Println("Body read error.")
	// }
	port := strings.Split(payload.Request.Host, ":")[1]
	t.Log.Println("Request host port: " + port)
	

	var uri []string
	if payload.Request != nil {
		uri = strings.Split(payload.Request.URL.RequestURI(), "/")
		abstractMessage.TranslatedMessage.Name = uri[len(uri)-1]
	} else {
		responseBody, err := httputil.DumpResponse(payload.Response, true)
		if err != nil {
			t.Log.Println("Could not dump payload request.")
		}

		t.Log.Println("Response payload: " + string(responseBody))
	}
	
	return *new(AbstractMessage)
}

func (t *GRPCTranslator) generateMessageId() uint64 {
	return t.index.Add(1)
}