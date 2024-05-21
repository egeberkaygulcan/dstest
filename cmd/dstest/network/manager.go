package network

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"log"
)

type Manager struct {
	Config       *config.Config
	Log          *log.Logger
	Interceptors []*Interceptor
}

func (nm *Manager) Init(config *config.Config) {
	numReplicas := config.ProcessConfig.NumReplicas

	nm.Config = config
	nm.Interceptors = make([]*Interceptor, numReplicas)

	// FIXME: This is a temporary solution to avoid nil pointer dereference
	nm.Log = log.New(log.Writer(), "[NetworkManager] ", log.LstdFlags)

	fmt.Printf("Config: %+v\n", config.NetworkConfig)
	nm.Log.Printf("Config: %+v\n", config.NetworkConfig)

	// create the interceptors
	for i := 0; i < numReplicas; i++ {
		nm.Interceptors[i] = new(Interceptor)
		nm.Interceptors[i].Init(i, nm.Config.NetworkConfig.BaseInterceptorPort+i, nm)
	}

	nm.Log.Println("Network manager initialized")
}

func (nm *Manager) Run() {
	// Run interceptors
	for i := 0; i < len(nm.Interceptors); i++ {
		go nm.Interceptors[i].Run()
	}

	nm.Log.Println("Network manager running")
}
