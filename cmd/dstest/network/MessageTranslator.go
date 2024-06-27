package network

import (
	// "io"
	"log"
	// "net/http/httputil"
	"os"
	"strings"
)

type MessageTranslator interface {
	Translate(*Message) *Message
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
	Log *log.Logger
}

func newGRPCTranslator() *GRPCTranslator {
	translator := new(GRPCTranslator)
	translator.Log = log.New(os.Stdout, "[GRPCTranslator]", log.LstdFlags)
	return translator
}

func (t *GRPCTranslator) Translate(message *Message) *Message {
	payload := message.Payload.(Http2CPayload)

	message.Type = GRPC

	// buf := new(strings.Builder)
	// _, err := io.Copy(buf, payload.Request.Body)
	// if err != nil {
	// 	t.Log.Println("Body read error.")
	// }

	var uri []string
	if payload.Request != nil {
		uri = strings.Split(payload.Request.URL.RequestURI(), "/")
		message.Name = uri[len(uri)-1]

		// requestBody, err := httputil.DumpRequest(payload.Request, true)
		// if err != nil {
		// 	t.Log.Println("Could not dump payload request.")
		// }

		// t.Log.Printf("Message sender: %d", message.Sender)
		// t.Log.Printf("Message receiver: %d", message.Receiver)
		// // t.Log.Println("Request payload: " + string(requestBody))

		// body, _ := io.ReadAll(payload.Request.Body)
		// t.Log.Println("Request body:")
		// for i := 0; i < len(body); i++ {
		// 	t.Log.Printf("%08b\n", body[i])
		// }
		// t.Log.Println("---------------")
	}
	// } else {
	// 	responseBody, err := httputil.DumpResponse(payload.Response, true)
	// 	if err != nil {
	// 		t.Log.Println("Could not dump payload request.")
	// 	}

	// 	t.Log.Println("Response payload: " + string(responseBody))
	// }
	
	return message
}