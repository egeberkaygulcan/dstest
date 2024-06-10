package network

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io"
	"io/ioutil"
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
	Request *http.Request
	Writer  http.ResponseWriter
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
	mux.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
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
		body, err := ioutil.ReadAll(request.Body)
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

		/**
		 * This is where we should stop execution and queue the message in the network manager
		 */
		hi.NetworkManager.Router.QueueMessage(Message{
			Sender:   -1,
			Receiver: thisNodePort - 6000,
			Payload:  Http2CPayload{Request: request, Writer: w},
		})

		// send the request
		resp, err := client.Do(proxyRequest)
		if err != nil {
			hi.Log.Fatalf("Error sending request to actual node: %s\n", err)
		}

		// send the response to the writer
		body = make([]byte, 1024*1024) // FIXME: make this dynamic
		_, err = resp.Body.Read(body)

		if err != nil {
			hi.Log.Fatalf("Error reading response: %s\n", err)
		}
		//hi.Log.Printf("Response: %s\n", body)

		w.WriteHeader(resp.StatusCode)
		_, err = w.Write(body)
		if err != nil {
			return
		}
		// close the connection
		defer resp.Body.Close()
	})

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
