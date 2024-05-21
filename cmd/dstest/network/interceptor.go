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
	Port           int
	ID             int
}

func (ni *Interceptor) Init(id int, port int, nm *Manager) {
	ni.ID = id
	ni.Port = port
	ni.NetworkManager = nm
	logHeader := fmt.Sprintf("[Interceptor %d] ", id)
	ni.Log = log.New(log.Writer(), logHeader, log.LstdFlags)
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
}
