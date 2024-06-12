package network

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

/**
 * Http2CInterceptor is an HTTP/2 Cleartext Interceptor
 * It is used to intercept HTTP/2 requests over cleartext (no TLS)
 * This is used e.g. by Ratis
 */
type Http2CInterceptor struct {
	BaseInterceptor
	Server *http.Server
}

type Http2CPayload struct {
	Request  *http.Request
	Response *http.Response
	Writer   http.ResponseWriter
}

// Check if BaseInterceptor implements Interceptor interface
var _ Interceptor = (*Http2CInterceptor)(nil)

func (hi *Http2CInterceptor) Init(id int, port int, nm *Manager) {
	logPrefix := fmt.Sprintf("[HTTP2C Interceptor %d] ", id)
	logger := log.New(log.Writer(), logPrefix, log.LstdFlags)
	hi.BaseInterceptor.Init(id, port, nm, logger)

	// create the multiplexer
	mux := http.NewServeMux()

	// handle all requests
	mux.HandleFunc("/", http2CRequestHandler(hi))

	h2s := &http2.Server{}

	// create the server
	hi.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        h2c.NewHandler(mux, h2s),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (hi *Http2CInterceptor) Run() (err error) {
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

func http2CRequestHandler(hi *Http2CInterceptor) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		hi.Log.Printf("Received request from %s: %s\n", request.RemoteAddr, request.URL.Path)

		// create connection to the node we're MITMing
		thisNodePort := hi.ID + 6000
		thisNodeAddr := fmt.Sprintf("localhost:%d", thisNodePort)
		//hi.Log.Printf("Connecting to actual node: %s\n", thisNodeAddr)
		client := http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctc context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctc, "tcp", thisNodeAddr)
				},
			},
		}

		// we need to buffer the body if we want to read it here and send it in the request
		body, err := io.ReadAll(request.Body)
		if err != nil {
			hi.Log.Fatalf("Error reading request body: %s\n", err)
		}

		// you can reassign the body if you need to parse it as multipart
		request.Body = io.NopCloser(bytes.NewReader(body))

		// create a new url from the raw RequestURI sent by the client
		url := fmt.Sprintf("http://localhost:%d%s", thisNodePort, request.RequestURI)

		proxyRequest, err := http.NewRequest(request.Method, url, bytes.NewReader(body))

		// we may want to filter some headers, otherwise we could just use a shallow copy
		// proxyRequest.Header = request.Header
		proxyRequest.Header = make(http.Header)
		for h, val := range request.Header {
			proxyRequest.Header[h] = val
		}

		// queue the request in the network manager
		awaitSendRequest := make(chan struct{})
		hi.NetworkManager.Router.QueueMessage(Message{
			Sender:   -1,
			Receiver: thisNodePort - 6000,
			Payload:  Http2CPayload{Request: proxyRequest, Writer: w, Response: nil},
			Send:     awaitSendRequest,
		})
		//<-awaitSendRequest

		// send the request
		resp, err := client.Do(proxyRequest)
		if err != nil {
			hi.Log.Fatalf("Error sending request to actual node: %s\n", err)
		}

		// create a buffer to hold the response body
		var buffer bytes.Buffer

		// copy the response body to the buffer
		_, err = io.Copy(&buffer, resp.Body)
		if err != nil {
			hi.Log.Fatalf("Error reading response: %s\n", err)
		}

		// convert the buffer to a byte slice
		body = buffer.Bytes()

		if err != nil {
			hi.Log.Fatalf("Error reading response: %s\n", err)
		}

		// queue sending the response in the network manager
		awaitSendResponse := make(chan struct{})
		hi.NetworkManager.Router.QueueMessage(Message{
			Sender:   -1,
			Receiver: thisNodePort - 6000,
			Payload:  Http2CPayload{Response: resp, Writer: w, Request: nil},
			Send:     awaitSendResponse,
		})
		//<-awaitSendResponse

		// send the response
		w.WriteHeader(resp.StatusCode)
		_, err = w.Write(body)
		if err != nil {
			return
		}

		// close the connection
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				hi.Log.Fatalf("Error closing response body: %s\n", err)
			}
		}(resp.Body)
	}
}

func dummyScheduler() {
	// dummy function to satisfy the interface
	fmt.Printf("Dummy scheduler invoked!!\n")
}
