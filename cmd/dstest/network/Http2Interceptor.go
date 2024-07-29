package network

// import (
// 	// "bufio"
// 	"bytes"
// 	"context"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net"
// 	"os"
// 	"sync"
// 	"time"

// 	"golang.org/x/net/http2"
// )

// type Http2Interceptor struct {
// 	BaseInterceptor
// 	Listener *net.TCPListener
// 	Ctx      context.Context
// 	Cancel   context.CancelFunc
// 	Conn   	 net.Conn
// 	Close    bool
// 	// WG       sync.WaitGroup
// }

// // Check if HttpInterceptor implements Interceptor interface
// var _ Interceptor = (*Http2Interceptor)(nil)

// func (hi *Http2Interceptor) Init(id int, port int, nm *Manager) {
// 	logPrefix := fmt.Sprintf("[HTTP Interceptor %d] ", id)
// 	logger := log.New(log.Writer(), logPrefix, log.LstdFlags)
// 	hi.BaseInterceptor.Init(id, port, nm, logger)
// 	hi.Conn = nil
// 	hi.Close = false
// 	hi.Ctx, hi.Cancel = context.WithCancel(context.Background())
// }

// func (hi *Http2Interceptor) Run() error {
// 	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", hi.Port))
// 	var err error
// 	hi.Listener, err = net.ListenTCP("tcp", addr)
// 	if err != nil {
// 		hi.Log.Println("Error while listening to TCP: ", err)
// 		return err
// 	}

// 	go func(ctx context.Context) {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				if hi.Listener != nil {
// 					conn, err := hi.Listener.Accept()
// 					// hi.Log.Println("Accepting new connection")
// 					if err != nil {
// 						if errors.Is(err, io.EOF) { // errors.Is(err, net.ErrClosed)
// 							hi.Log.Println("Error while accepting TCP connection: ", err)
// 						}
// 						continue
// 					}

// 					// hi.WG.Add(1)
// 					go func() {
// 						hi.handleConn(conn.(*net.TCPConn))
// 						// hi.WG.Done()
// 					}()
// 				}
// 			}
// 		}
// 	}(hi.Ctx)

// 	return nil
// }

// func (hi *Http2Interceptor) handleConn(conn *net.TCPConn) {
// 	defer func() {
// 		_ = conn.Close()
// 	}()

// 	pair := hi.NetworkManager.PortMap[hi.Port]
// 	f := http2.NewFramer(conn, conn)
// 	messageCount := 0
// 	for !hi.Close {
// 		// hi.Log.Println("Handling Request:")
// 		// buffer := make([]byte, 4096)
// 		// length, err := conn.Read(buffer)
// 		// if err != nil {
// 		// 	continue
// 		// }

// 		// f := http2.NewFramer(conn, conn)

// 		if messageCount == 0 {
// 			const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
// 			b := make([]byte, len(preface))
// 			if _, err := io.ReadFull(conn, b); err != nil {
// 				continue
// 			}

// 			if string(b) != preface {
// 				hi.Log.Println("Invalid preface")
// 				continue
// 			} else {
// 				hi.Log.Println("Valid preface")
// 			}

// 			err := f.WriteSettings()
// 			if err != nil {
// 				hi.Log.Printf("Error while writing HTTP2 settings: %s\n", err)
// 			}
			
// 			messageCount++
// 		}

// 		frames := make([]http2.Frame, 0)
// 		for {
// 			frame, err := f.ReadFrame()
// 			if err != nil {
// 				break
// 			}
// 			switch frame.Header().Type {
// 			case http2.FrameGoAway:
// 				return
// 			case http2.FrameSettings:
// 				err := f.WriteSettingsAck()
// 				if err != nil {
// 					hi.Log.Printf("Error while writing HTTP2 settings: %s\n", err)
// 				}
// 			case http2.FrameHeaders:
// 				hi.Log.Println("---------------------------------")

// 				if len(frames) > 0 {
// 					hi.NetworkManager.Router.QueueMessage(&Message{
// 						Sender:    pair.Sender,
// 						Receiver:  pair.Receiver,
// 						Payload:   frames,
// 						Type:      "",
// 						Name:      "",
// 						MessageId: hi.NetworkManager.GenerateUniqueId(),
// 					})
// 					frames = make([]http2.Frame, 0)
// 				} 
// 				frames = append(frames, frame)
// 			case http2.FrameData:
// 				frames = append(frames, frame)
// 			}
// 			hi.Log.Println(frame.Header().Type)
// 		}

// 		// str := string(buffer[:length])
// 		// hi.Log.Println(conn.RemoteAddr().String())
//         // hi.Log.Printf("Received command %d\t:%s\n", length, str)

// 		// pair := hi.NetworkManager.PortMap[hi.Port]
// 		// // queue the request in the network manager
// 		// hi.NetworkManager.Router.QueueMessage(&Message{
// 		// 	Sender:    pair.Sender,
// 		// 	Receiver:  pair.Receiver,
// 		// 	Payload:   buffer[:length],
// 		// 	Type:      "",
// 		// 	Name:      "",
// 		// 	MessageId: hi.NetworkManager.GenerateUniqueId(),
// 		// })

// 		// const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
// 		// b := make([]byte, len(preface))
// 		// if _, err := io.ReadFull(conn, b); err != nil {
// 		// 	hi.Log.Println(err)
// 		// }

// 		// if string(b) != preface {
// 		// 	hi.Log.Println("Invalid preface")
// 		// } else {
// 		// 	hi.Log.Println("Valid preface")
// 		// }

// 		// framer := http2.NewFramer(conn, conn)
// 		// frame, err := framer.ReadFrame()
// 		// if err != nil {
// 		// 	hi.Log.Println(frame, err)

// 		// 	for err == nil {
// 		// 		frame, err = framer.ReadFrame()
// 		// 		hi.Log.Println(frame, err)
// 		// 	}
// 		// }

// 	}
// }

// func (hi *Http2Interceptor) SendMessage(message *Message) {
// 	if hi.Conn == nil {
// 		thisNodePort := message.Receiver + hi.NetworkManager.Config.NetworkConfig.BaseReplicaPort
// 		var err error
// 		hi.Conn, err = net.Dial("tcp", fmt.Sprintf(":%d", thisNodePort))
// 		if err != nil {
// 			hi.Log.Println("Error while establishing connection with the original node: ", err)
// 		}
// 	}

// 	// framer := http2.NewFramer(hi.Conn, hi.Conn)
// 	frames := message.Payload.([]http2.Frame)

// 	for _, frame := range frames {
// 		switch frame.Header().Type {
// 		case http2.FrameHeaders:
// 			frame = frame.(*http2.HeadersFrame)
// 			// framer.WriteHeaders() // TODO - New HeadersFrameParam
// 		}
// 	}

// 	// hi.Log.Printf("Sending message: %d - %d", message.Sender, message.Receiver)
// 	// _, err := io.Copy(hi.Conn, bytes.NewReader(message.Payload))
// 	// if err != nil {
// 	// 	hi.Log.Println("Error while sending message: ", err)
// 	// 	str := string(message.Payload)
// 	// 	hi.Log.Println("Message: ", str)
// 	// }
// }

// func (hi *Http2Interceptor) Shutdown() {
// 	hi.Cancel()
// 	hi.Listener.Close()
// 	hi.Close = true
// 	// hi.WG.Wait()
// }

// func (hi *Http2Interceptor) handleHttp2(conn net.Conn) error {
// 	defer func() {
// 		_ = conn.Close()
// 	}()

// 	dataBuffer := bytes.NewBuffer(make([]byte, 0))
// 	reader := io.TeeReader(conn, dataBuffer)

// 	f := http2.NewFramer(conn, conn)

// 	err := f.WriteSettings()
// 	if err != nil {
// 		hi.Log.Printf("Error while writing HTTP2 settings: %s\n", err)
// 		return err
// 	}

// 	// err = f.WriteSettingsAck()
// 	// if err != nil {
// 	// 	hi.Log.Printf("Error while writing HTTP2 settings Ack: %s\n", err)
// 	// 	return err
// 	// }

// 	f = http2.NewFramer(io.Discard, reader)

// 	// Intercept and wait
// 	pair := hi.NetworkManager.PortMap[hi.Port]
// 	thisNodePort := pair.Receiver + hi.NetworkManager.Config.NetworkConfig.BaseReplicaPort
// 	// queue the request in the network manager
// 	awaitSendRequest := make(chan struct{})
// 	hi.NetworkManager.Router.QueueMessage(&Message{
// 		Sender:    pair.Sender,
// 		Receiver:  pair.Receiver,
// 		// Payload:   f,
// 		Type:      "",
// 		Name:      "",
// 		MessageId: hi.NetworkManager.GenerateUniqueId(),
// 		// Send:      awaitSendRequest,
// 	})
// 	<-awaitSendRequest

// 	hi.Log.Printf("Sending from: %d", pair.Sender)
// 	dialer, err := net.Dial("tcp", fmt.Sprintf(":%d", thisNodePort))
// 	if err != nil {
// 		hi.Log.Printf("Error while sending to original node: %s\n", err)
// 		return err
// 	}

// 	_ = dialer.SetReadDeadline(time.Now().Add(time.Duration(hi.NetworkManager.Config.TestConfig.WaitDuration) * time.Millisecond))
// 	_ = dialer.SetWriteDeadline(time.Now().Add(time.Duration(hi.NetworkManager.Config.TestConfig.WaitDuration) * time.Millisecond))

// 	wg := sync.WaitGroup{}
// 	wg.Add(1)
// 	dataSent := int64(0)
// 	go func(dataSent *int64) {
// 		*dataSent, err = io.Copy(conn, dialer)
// 		wg.Done()
// 	}(&dataSent)

// 	_, err = io.Copy(dialer, io.MultiReader(dataBuffer, conn))
// 	hi.Log.Println("Copy 2 completed.")
// 	wg.Wait()
// 	hi.Log.Println("WG, waited")

// 	if errors.Is(err, os.ErrDeadlineExceeded) && dataSent > 0 {
// 		return nil
// 	}

// 	return err
// }