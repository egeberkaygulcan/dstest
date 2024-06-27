package network

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	// "github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling"
)

type SenderReceiverPair struct {
	Sender int
	Receiver int
}

type Manager struct {
	Config        *config.Config
	Log           *log.Logger
	Router        *Router
	Interceptors  []Interceptor
	MessageQueues []*MessageQueue
	index atomic.Uint64
	WaitGroup	  sync.WaitGroup
	ReplicaIds []int
	PortMap map[int]SenderReceiverPair

	// Scheduler     scheduling.Scheduler
}

func (nm *Manager) Init(config *config.Config, replicaIds []int) {
	numReplicas := config.ProcessConfig.NumReplicas

	nm.Config = config
	nm.Router = new(Router)
	nm.Interceptors = make([]Interceptor, numReplicas * (numReplicas - 1))
	nm.MessageQueues = make([]*MessageQueue, numReplicas)
	nm.ReplicaIds = replicaIds
	// nm.Scheduler = new(scheduling.BasicScheduler)

	nm.Router.Init(nm, numReplicas)

	// create the interceptors and message queues
	nm.PortMap = make(map[int]SenderReceiverPair)
	for i := 0; i < numReplicas; i++ {
		nm.MessageQueues[i] = new(MessageQueue)
		nm.MessageQueues[i].Init()
		for j := 0; j < numReplicas; j++ {
			if i != j {
				id := i*numReplicas+j
				nm.Interceptors[id] = new(Http2CInterceptor)
				nm.Interceptors[id].Init(id, nm.Config.NetworkConfig.BaseInterceptorPort+id, nm)
				nm.PortMap[nm.Config.NetworkConfig.BaseInterceptorPort+id] = SenderReceiverPair{Sender: i, Receiver: j}
			}
		}
	}

	// FIXME: This is a temporary solution to avoid nil pointer dereference
	nm.Log = log.New(log.Writer(), "[NetworkManager] ", log.LstdFlags)

	nm.Log.Println("Network manager initialized")
}

func (nm *Manager) Run() {
	// Run interceptors
	for i := 0; i < len(nm.Interceptors); i++ {
		nm.WaitGroup.Add(1)
		go func(index int) {
			nm.Interceptors[index].Run()
			nm.WaitGroup.Done()
		}(i)
	}

	nm.Log.Println("Network manager running")
	nm.WaitGroup.Wait()
}

func (nm *Manager) Shutdown() {
	for _, interceptor := range nm.Interceptors {
		interceptor.Shutdown()
	}
}

func (nm *Manager) GenerateUniqueId() uint64 {
	return nm.index.Add(1) 
}

func (nm *Manager) SendMessage(messageId uint64) {
	for _, mq := range nm.MessageQueues {
		if mq.Peek() != nil {
			if mq.Peek().MessageId == messageId {
				message := mq.PopFront()
				message.SendMessage()
			}
		}
	}
}

func (nm *Manager) GetActions() []*Message {
	var actions []*Message

	delayMessage := &(Message{
		Sender:   -1,
		Receiver: -1,
		Payload:  Http2CPayload{Request: nil, Writer: nil, Response: nil},
		Type: "Delay",
		Name: "Delay",
		MessageId: nm.GenerateUniqueId(),
		Send:     nil,
	})
	actions = append(actions, delayMessage)

	for _, mq := range nm.MessageQueues {
		action := mq.Peek()

		if action != nil {
			actions = append(actions, action)
		}
	}

	return actions
}