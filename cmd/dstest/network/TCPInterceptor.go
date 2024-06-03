package network

import (
	"fmt"
	"log"
	"net"
)

type TCPInterceptor struct {
	BaseInterceptor
}

// Check if BaseInterceptor implements Interceptor interface
var _ Interceptor = (*TCPInterceptor)(nil)

func (ni *TCPInterceptor) Init(id int, port int, nm *Manager) {
	logPrefix := fmt.Sprintf("[TCP Interceptor %d] ", id)
	logger := log.New(log.Writer(), logPrefix, log.LstdFlags)
	ni.BaseInterceptor.Init(id, port, nm, logger)
}

func (ni *TCPInterceptor) Run() (err error) {
	err = ni.BaseInterceptor.Run()
	if err != nil {
		return err
	}

	// log the port
	ni.Log.Printf("Running TCP interceptor on port %d\n", ni.Port)

	// Start listening on the port
	portSpecification := fmt.Sprintf(":%d", ni.Port)
	listener, err := net.Listen("tcp", portSpecification)

	// Check for errors
	if err != nil {
		ni.Log.Fatalf("Error listening on port %d: %s\n", ni.Port, err.Error())
		return err
	}

	ni.Log.Printf("Listening on port %d\n", ni.Port)

	// Close the listener when the function returns
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			ni.Log.Fatalf("Error closing listener on port %d: %s\n", ni.Port, err.Error())
		}
	}(listener)

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			ni.Log.Fatalf("Error accepting connection on port %d: %s\n", ni.Port, err.Error())
		}

		go ni.handleConnection(conn)
	}
}

func (ni *TCPInterceptor) handleConnection(conn net.Conn) {
	//ni.Log.Printf("Handling connection from %s\n", conn.RemoteAddr().String())

	// echo back the message
	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		ni.Log.Fatalf("Error reading from connection: %s\n", err.Error())
	}

	ni.Log.Printf("Received from %s: %s\n", conn.RemoteAddr().String(), string(buf[:n]))

	// Push the message to the receiver's message queue
	ni.NetworkManager.Router.QueueMessage(Message{
		Sender:   conn.RemoteAddr().(*net.TCPAddr).Port - ni.NetworkManager.Config.NetworkConfig.BaseReplicaPort,
		Receiver: ni.ID,
		Payload:  buf[:n],
	})

	/*
		// Send the message back to the sender
		_, err = conn.Write(buf[:n])
		if err != nil {
			ni.Log.Fatalf("Error writing to connection: %s\n", err.Error())
		}
		ni.Log.Printf("Sent message: %s\n", string(buf[:n]))
	*/

	err = conn.Close()
	if err != nil {
		ni.Log.Fatalf("Error closing connection: %s\n", err.Error())
	}

	//ni.Log.Printf("Connection closed\n")
}
