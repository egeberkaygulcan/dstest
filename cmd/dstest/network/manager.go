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
	MessageType MessageType
	VectorClocks map[int]map[int]int
}

func (nm *Manager) Init(config *config.Config, replicaIds []int) {
	numReplicas := config.ProcessConfig.NumReplicas

	nm.Config = config
	nm.MessageType = MessageType(config.NetworkConfig.MessageType)
	nm.Router = new(Router)
	nm.Interceptors = make([]Interceptor, numReplicas * (numReplicas - 1))
	nm.MessageQueues = make([]*MessageQueue, numReplicas)
	nm.ReplicaIds = replicaIds
	nm.VectorClocks = make(map[int]map[int]int)

	nm.Router.Init(nm, numReplicas)

	// create the interceptors and message queues
	nm.PortMap = make(map[int]SenderReceiverPair)
	k := 0
	for i := 0; i < numReplicas; i++ {
		nm.MessageQueues[i] = new(MessageQueue)
		nm.MessageQueues[i].Init()
		for j := 0; j < numReplicas; j++ {
			if i != j {
				id := i*numReplicas+j
				// nm.Interceptors[k] = new(Http2CInterceptor)
				nm.Interceptors[k] = new(HttpInterceptor)
				nm.Interceptors[k].Init(id, nm.Config.NetworkConfig.BaseInterceptorPort+id, nm)
				nm.PortMap[nm.Config.NetworkConfig.BaseInterceptorPort+id] = SenderReceiverPair{Sender: i, Receiver: j}
				k++
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

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (nm *Manager) updateVectorClocks(sender, receiver int) {
	// Update sender clock
	nm.VectorClocks[sender][sender]++

	// Update receiver clock
	for _, id := range nm.ReplicaIds {
		if id != receiver {
			nm.VectorClocks[receiver][id] = max(nm.VectorClocks[receiver][id], nm.VectorClocks[sender][id])
		} else {
			nm.VectorClocks[receiver][receiver]++
		}
	}
}

func (nm *Manager) SendMessage(messageId uint64) {
	for _, mq := range nm.MessageQueues {
		if mq.Peek() != nil {
			if mq.Peek().MessageId == messageId {
				message := mq.PopFront()
				message.SendMessage()
				nm.updateVectorClocks(message.Sender, message.Receiver)
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