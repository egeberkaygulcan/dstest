package network

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling"
	"log"
)

type Manager struct {
	Config        *config.Config
	Log           *log.Logger
	Router        *Router
	Interceptors  []Interceptor
	MessageQueues []*MessageQueue
	Scheduler     scheduling.Scheduler
}

func (nm *Manager) Init(config *config.Config) {
	numReplicas := config.ProcessConfig.NumReplicas

	nm.Config = config
	nm.Router = new(Router)
	nm.Interceptors = make([]Interceptor, numReplicas)
	nm.MessageQueues = make([]*MessageQueue, numReplicas)
	nm.Scheduler = new(scheduling.BasicScheduler)

	nm.Router.Init(nm, numReplicas)

	// create the interceptors and message queues
	for i := 0; i < numReplicas; i++ {
		nm.MessageQueues[i] = new(MessageQueue)
		nm.Interceptors[i] = new(Http2CInterceptor)
		nm.MessageQueues[i].Init()
		nm.Interceptors[i].Init(i, nm.Config.NetworkConfig.BaseInterceptorPort+i, nm)
	}

	// FIXME: This is a temporary solution to avoid nil pointer dereference
	nm.Log = log.New(log.Writer(), "[NetworkManager] ", log.LstdFlags)

	nm.Log.Println("Network manager initialized")
}

func (nm *Manager) Run() {
	// Run interceptors
	for i := 0; i < len(nm.Interceptors); i++ {
		go nm.Interceptors[i].Run()
	}

	nm.Log.Println("Network manager running")
}
