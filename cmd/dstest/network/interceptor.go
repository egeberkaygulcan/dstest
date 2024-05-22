package network

import (
	"fmt"
	"log"
	"net"
)

type Interceptor struct {
	isInitialized  bool
	Log            *log.Logger
	NetworkManager *Manager
	MessageQueue   *MessageQueue
	Port           int
	ID             int
}

func (ni *Interceptor) Init(id int, port int, nm *Manager) {
	ni.ID = id
	ni.Port = port
	ni.NetworkManager = nm
	ni.MessageQueue = nm.MessageQueues[id]
	logPrefix := fmt.Sprintf("[Interceptor %d] ", id)
	ni.Log = log.New(log.Writer(), logPrefix, log.LstdFlags)
	ni.isInitialized = true
}

func (ni *Interceptor) Run() {
	// check if the interceptor is initialized
	if !ni.isInitialized {
		ni.Log.Fatalf("Interceptor is not initialized\n")
		return
	}

	// log the port
	fmt.Printf("Running network interceptor on port %d\n", ni.Port)

	ni.Log.Printf("Running network interceptor on port %d\n", ni.Port)

	// Start listening on the port
	portSpecification := fmt.Sprintf(":%d", ni.Port)
	listener, err := net.Listen("tcp", portSpecification)

	// Check for errors
	if err != nil {
		ni.Log.Fatalf("Error listening on port %d: %s\n", ni.Port, err.Error())
	}

	// Close the listener when the function returns
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			ni.Log.Fatalf("Error closing listener on port %d: %s\n", ni.Port, err.Error())
		}
	}(listener)

	ni.Log.Printf("Listening on port %d\n", ni.Port)

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			ni.Log.Fatalf("Error accepting connection on port %d: %s\n", ni.Port, err.Error())
		}

		go ni.handleConnection(conn)
	}
}

func (ni *Interceptor) handleConnection(conn net.Conn) {
	ni.Log.Printf("Handling connection from %s\n", conn.RemoteAddr().String())

	ni.Log.Printf("conn:" + conn.LocalAddr().String())

	// echo back the message
	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		ni.Log.Fatalf("Error reading from connection: %s\n", err.Error())
	}

	message := Message{
		Sender:   conn.RemoteAddr().(*net.TCPAddr).Port - ni.NetworkManager.Config.NetworkConfig.BaseReplicaPort,
		Receiver: ni.ID,
		Payload:  buf[:n],
	}
	ni.MessageQueue.PushBack(message)
	ni.MessageQueue.Print(ni.Log)

	ni.Log.Printf("Received message: %s\n", string(buf[:n]))

	_, err = conn.Write(buf[:n])
	if err != nil {
		ni.Log.Fatalf("Error writing to connection: %s\n", err.Error())
	}

	ni.Log.Printf("Sent message: %s\n", string(buf[:n]))

	err = conn.Close()
	if err != nil {
		ni.Log.Fatalf("Error closing connection: %s\n", err.Error())
	}

	ni.Log.Printf("Connection closed\n")
}
