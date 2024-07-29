package network

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

type HttpInterceptor struct {
	BaseInterceptor
	Listener *net.TCPListener
	Ctx      context.Context
	Cancel   context.CancelFunc
	// WG       sync.WaitGroup
}

// Check if HttpInterceptor implements Interceptor interface
var _ Interceptor = (*HttpInterceptor)(nil)

func (hi *HttpInterceptor) Init(id int, port int, nm *Manager) {
	logPrefix := fmt.Sprintf("[HTTP Interceptor %d] ", id)
	logger := log.New(log.Writer(), logPrefix, log.LstdFlags)
	hi.BaseInterceptor.Init(id, port, nm, logger)
	hi.Ctx, hi.Cancel = context.WithCancel(context.Background())
}

func (hi *HttpInterceptor) Run() error {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", hi.Port))
	var err error
	hi.Listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		hi.Log.Printf("Error while listening to TCP: %s\n", err)
		return err
	}

	go func() {
		for {
			if hi.Listener != nil {
				conn, err := hi.Listener.Accept()
				// hi.Log.Println("Accepting new connection")
				if err != nil {
					if errors.Is(err, io.EOF) {//|| errors.Is(err, net.ErrClosed) {
						hi.Log.Printf("Error while accepting TCP connection: %s\n", err)
					}
					continue
				}

				go hi.handleConn(conn.(*net.TCPConn))
				
			}
		}
	}()

	return nil
}

func (hi *HttpInterceptor) handleConn(conn *net.TCPConn) {
	defer func() {
		_ = conn.Close()
	}()
	var err error
	const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"


	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		hi.Log.Println("Error while reading payload: ", err)
	}

	b := buf[:len(preface)]
	payload := buf[:n]

	if string(b) == preface || hi.NetworkManager.Config.NetworkConfig.MessageType == string(GRPC) {
		err = hi.handleHttp2(bytes.NewBuffer(payload), conn)
		if err != nil {
			hi.Log.Println("Error while handling http2: ", err)
		}
	} else {
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(payload)))
		if err != nil {
			hi.Log.Printf("Error while creating new read request: %s\n", err)
			return
		}

		err = hi.handleHttpReq(req, conn)
		hi.Log.Println("Error while handling http: ", err)
	}
}

func (hi *HttpInterceptor) handleHttpReq(req *http.Request, w net.Conn) error {
	// resp := &http.Response{}
	// // Intercept
	// // TODO
	// // Required to forward the request
	// req.RequestURI = ""
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return err
	// }
	// return resp.Write(w)
	return nil
}

func (hi *HttpInterceptor) handleHttp2(initial io.Reader, conn net.Conn) error {
	defer func() {
		_ = conn.Close()
	}()

	dataBuffer := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(conn, dataBuffer)

	f := http2.NewFramer(conn, conn)
	err := f.WriteSettings()
	if err != nil {
		hi.Log.Printf("Error while writing HTTP2 settings: %s\n", err)
		return err
	}

	// err = f.WriteSettingsAck()
	// if err != nil {
	// 	hi.Log.Printf("Error while writing HTTP2 settings Ack: %s\n", err)
	// 	return err
	// }

	f = http2.NewFramer(io.Discard, reader)

	// Intercept and wait
	pair := hi.NetworkManager.PortMap[hi.Port]
	thisNodePort := pair.Receiver + hi.NetworkManager.Config.NetworkConfig.BaseReplicaPort
	// queue the request in the network manager
	awaitSendRequest := make(chan struct{})
	hi.NetworkManager.Router.QueueMessage(&Message{
		Sender:    pair.Sender,
		Receiver:  pair.Receiver,
		Payload:   f,
		Type:      "",
		Name:      "",
		MessageId: hi.NetworkManager.GenerateUniqueId(),
		Send:      awaitSendRequest,
	})
	<-awaitSendRequest

	hi.Log.Printf("Sending from: %d", pair.Sender)
	dialer, err := net.Dial("tcp", fmt.Sprintf(":%d", thisNodePort))
	if err != nil {
		hi.Log.Printf("Error while sending to original node: %s\n", err)
		return err
	}

	_ = dialer.SetReadDeadline(time.Now().Add(time.Duration(hi.NetworkManager.Config.TestConfig.WaitDuration) * time.Millisecond))
	_ = dialer.SetWriteDeadline(time.Now().Add(time.Duration(hi.NetworkManager.Config.TestConfig.WaitDuration) * time.Millisecond))

	wg := sync.WaitGroup{}
	wg.Add(1)
	dataSent := int64(0)
	go func(dataSent *int64) {
		*dataSent, err = io.Copy(conn, dialer)
		wg.Done()
	}(&dataSent)

	_, err = io.Copy(dialer, io.MultiReader(initial, dataBuffer, conn))
	hi.Log.Println("Copy 2 completed.")
	wg.Wait()
	hi.Log.Println("WG, waited")

	if errors.Is(err, os.ErrDeadlineExceeded) && dataSent > 0 {
		return nil
	}

	return err
}

func (hi *HttpInterceptor) Shutdown() {
	hi.Cancel()
	hi.Listener.Close()
	// hi.WG.Wait()
}

// package network

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"
// )

// type HttpAInterceptor struct {
// 	BaseInterceptor
// 	Server *http.Server
// }

// type HttpPayload struct {
// 	Request *http.Request
// 	Writer  http.ResponseWriter
// }

// // Check if BaseInterceptor implements Interceptor interface
// var _ Interceptor = (*HttpAInterceptor)(nil)

// func (hi *HttpAInterceptor) Init(id int, port int, nm *Manager) {
// 	logPrefix := fmt.Sprintf("[HTTP Interceptor %d] ", id)
// 	logger := log.New(log.Writer(), logPrefix, log.LstdFlags)
// 	hi.BaseInterceptor.Init(id, port, nm, logger)

// 	// create the multiplexer
// 	mux := http.NewServeMux()

// 	// handle all requests
// 	mux.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
// 		hi.Log.Printf("Received request: %s\n", request.URL.Path)
// 		// log the sender ip
// 		hi.Log.Printf("Sender IP: %s\n", request.RemoteAddr)

// 		remotePort, _ := strconv.Atoi(strings.Split(request.RemoteAddr, ":")[1])

// 		hi.NetworkManager.Router.QueueMessage(&Message{
// 			Sender:   remotePort,
// 			Receiver: hi.ID,
// 			Payload:  HttpPayload{Request: request, Writer: w},
// 			Type: "",
// 			Name: "",
// 			MessageId: hi.NetworkManager.GenerateUniqueId(),
// 		})
// 	})

// 	// create the server
// 	hi.Server = &http.Server{
// 		Addr:    fmt.Sprintf(":%d", port),
// 		Handler: mux,
// 	}
// }

// func (hi *HttpAInterceptor) Run() (err error) {
// 	// check if the interceptor is initialized
// 	if !hi.isInitialized {
// 		hi.Log.Fatalf("BaseInterceptor is not initialized\n")
// 		return
// 	}

// 	// log the port
// 	hi.Log.Printf("Running HTTP interceptor on port %d\n", hi.Port)

// 	go func() {
// 		err := hi.Server.ListenAndServe()
// 		if err != nil {
// 			hi.Log.Fatalf("Error listening on port %d: %s\n", hi.Port, err.Error())
// 		}
// 	}()

// 	return nil
// }

// func (hi *HttpAInterceptor) Shutdown() {

// }
