package network

import (
	"errors"
	"log"
)

type Interceptor interface {
	Init(id int, port int, nm *Manager)
	Run() (err error)
	Shutdown()
}

type BaseInterceptor struct {
	isInitialized  bool
	Log            *log.Logger
	NetworkManager *Manager
	Port           int
	ID             int
}

func (ni *BaseInterceptor) Init(id int, port int, nm *Manager, log *log.Logger) {
	ni.ID = id
	ni.Port = port
	ni.NetworkManager = nm
	ni.Log = log
	ni.isInitialized = true
}

func (ni *BaseInterceptor) Run() (err error) {
	// check if the interceptor is initialized
	if !ni.isInitialized {
		return errors.New("Interceptor is not initialized")
	}

	return nil
}
