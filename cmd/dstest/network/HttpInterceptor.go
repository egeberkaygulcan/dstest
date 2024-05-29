package network

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HttpInterceptor struct {
	BaseInterceptor
	Server *http.Server
}

type HttpPayload struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

// Check if BaseInterceptor implements Interceptor interface
var _ Interceptor = (*HttpInterceptor)(nil)

func (hi *HttpInterceptor) Init(id int, port int, nm *Manager) {
	logPrefix := fmt.Sprintf("[HTTP Interceptor %d] ", id)
	logger := log.New(log.Writer(), logPrefix, log.LstdFlags)
	hi.BaseInterceptor.Init(id, port, nm, logger)

	// create the multiplexer
	mux := http.NewServeMux()

	// handle all requests
	mux.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		hi.Log.Printf("Received request: %s\n", request.URL.Path)
		// log the sender ip
		hi.Log.Printf("Sender IP: %s\n", request.RemoteAddr)

		remotePort, _ := strconv.Atoi(strings.Split(request.RemoteAddr, ":")[1])

		hi.NetworkManager.Router.QueueMessage(Message{
			Sender:   remotePort,
			Receiver: hi.ID,
			Payload:  HttpPayload{Request: request, Writer: w},
		})
	})

	// create the server
	hi.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}

func (hi *HttpInterceptor) Run() (err error) {
	// check if the interceptor is initialized
	if !hi.isInitialized {
		hi.Log.Fatalf("BaseInterceptor is not initialized\n")
		return
	}

	// log the port
	hi.Log.Printf("Running HTTP interceptor on port %d\n", hi.Port)

	go func() {
		err := hi.Server.ListenAndServe()
		if err != nil {
			hi.Log.Fatalf("Error listening on port %d: %s\n", hi.Port, err.Error())
		}
	}()

	return nil
}
