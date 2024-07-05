package network

import (
	"fmt"
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

type Event struct {
	Prev *Event
	Sender int
	Receiver int
	MessageId uint64
	Name string
}

type Manager struct {
	Config        *config.Config
	Log           *log.Logger
	Router        *Router
	Interceptors  []Interceptor
	MessageQueues []*MessageQueue
	Index atomic.Uint64
	WaitGroup	  sync.WaitGroup
	ReplicaIds []int
	PortMap map[int]SenderReceiverPair
	MessageType MessageType
	// VectorClocks map[int]map[int]int
	ChainClocks [][]Event
}

func (nm *Manager) Init(config *config.Config, replicaIds []int) error {

	numReplicas := config.ProcessConfig.NumReplicas

	nm.Config = config
	nm.MessageType = MessageType(config.NetworkConfig.MessageType)
	nm.Router = new(Router)
	nm.Interceptors = make([]Interceptor, numReplicas * (numReplicas - 1))
	nm.MessageQueues = make([]*MessageQueue, numReplicas)
	nm.ReplicaIds = replicaIds
	// nm.VectorClocks = make(map[int]map[int]int)
	nm.ChainClocks = make([][]Event, 0)
	nm.Index.Store(0)

	nm.Router.Init(nm, numReplicas)

	// create the interceptors and message queues
	nm.PortMap = make(map[int]SenderReceiverPair)
	k := 0
	for i := 0; i < numReplicas; i++ {
		nm.MessageQueues[i] = new(MessageQueue)
		var err error = nil

		nm.MessageQueues[i].Init()
		for j := 0; j < numReplicas; j++ {
			if i != j {
				id := i*numReplicas + j
				nm.Interceptors[k], err = createInterceptor(config.NetworkConfig.Protocol)
				if err != nil {
					return fmt.Errorf("Error creating interceptor: %s", err.Error())
				}
				nm.Interceptors[k].Init(id, nm.Config.NetworkConfig.BaseInterceptorPort+id, nm)
				nm.PortMap[nm.Config.NetworkConfig.BaseInterceptorPort+id] = SenderReceiverPair{Sender: i, Receiver: j}
				k++
			}
		}
	}

	// FIXME: This is a temporary solution to avoid nil pointer dereference
	nm.Log = log.New(log.Writer(), "[NetworkManager] ", log.LstdFlags)

	nm.Log.Println("Network manager initialized")
	return nil
}

func (nm *Manager) Run() {
	// Run interceptors
	for i := 0; i < len(nm.Interceptors); i++ {
		nm.WaitGroup.Add(1)
		go func(index int) {
			err := nm.Interceptors[index].Run()
			if err != nil {
				errStr := fmt.Errorf("Error running interceptor: %s", err.Error())
				fmt.Println(errStr)
				return
			}
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

	for _, chain := range nm.ChainClocks {
		nm.Log.Println("Chain: ")
		for _, elem := range chain {
			nm.Log.Printf("(%d)-(%d)-(%s)-(%d)", elem.Sender, elem.Receiver, elem.Name, int(elem.MessageId))
		}
	}
}

func (nm *Manager) GenerateUniqueId() uint64 {
	return nm.Index.Add(1)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// func (nm *Manager) updateVectorClocks(sender, receiver int) {
// 	// Update sender clock
// 	nm.VectorClocks[sender][sender]++

// 	// Update receiver clock
// 	for _, id := range nm.ReplicaIds {
// 		if id != receiver {
// 			nm.VectorClocks[receiver][id] = max(nm.VectorClocks[receiver][id], nm.VectorClocks[sender][id])
// 		} else {
// 			nm.VectorClocks[receiver][receiver]++
// 		}
// 	}
// }

func (nm *Manager) UpdateChainClocks(sender, receiver int, messageId uint64, name string) {
	for i, chain := range nm.ChainClocks {
		if chain[len(chain)-1].Receiver == sender  {
				nm.ChainClocks[i] = append(nm.ChainClocks[i], Event{&chain[len(chain)-1], sender, receiver, messageId, name})
				return
		}
	}

	elem := []Event{Event{nil, sender, receiver, messageId, name}}
	nm.ChainClocks = append(nm.ChainClocks, elem)
}

func (nm *Manager) SendMessage(messageId uint64) {
	for _, mq := range nm.MessageQueues {
		if mq.Peek() != nil {
			if mq.Peek().MessageId == messageId {
				message := mq.PopFront()
				message.SendMessage()
				// nm.updateVectorClocks(message.Sender, message.Receiver)
			}
		}
	}
}

func (nm *Manager) GetActions() []*Message {
	var actions []*Message

	// delayMessage := &(Message{
	// 	Sender:   -1,
	// 	Receiver: -1,
	// 	Payload:  Http2CPayload{Request: nil, Writer: nil, Response: nil},
	// 	Type: "Delay",
	// 	Name: "Delay",
	// 	MessageId: uint64(0),
	// 	Send:     nil,
	// })
	// actions = append(actions, delayMessage)

	for _, mq := range nm.MessageQueues {
		action := mq.Peek()

		if action != nil {
			actions = append(actions, action)
		}
	}

	return actions
}
