package network

import (
	"fmt"
	"log"
)

type Router struct {
	RoutingTable   [][]bool
	NetworkManager *Manager
	MessageTranslator MessageTranslator
	Log            *log.Logger
}

func (r *Router) Init(NetworkManager *Manager, numReplicas int) {
	r.NetworkManager = NetworkManager
	r.RoutingTable = make([][]bool, numReplicas)
	r.MessageTranslator = NewMessageTranslator(GRPC)
	r.Log = log.New(log.Writer(), "[Router] ", log.LstdFlags)

	// create N*N routing table
	for i := 0; i < numReplicas; i++ {
		r.RoutingTable[i] = make([]bool, numReplicas)
		// initialize to True
		for j := 0; j < numReplicas; j++ {
			r.RoutingTable[i][j] = true
		}
	}

	//r.CreatePartitions([][]int{[]int{0, 1}, []int{2}})
	//r.PrintRoutingTable()
}

// queue message
func (r *Router) QueueMessage(m *Message) {
	// check if there is connectivity between sender and receiver
	tm := r.MessageTranslator.Translate(m)
	if r.HasConnectivity(tm.Sender, tm.Receiver) {
		r.NetworkManager.MessageQueues[tm.Receiver].PushBack(tm)
		r.Log.Printf("Queued message #%d from %d to %d: %s\n", r.NetworkManager.MessageQueues[tm.Receiver].Len(), tm.Sender, tm.Receiver, (tm.Payload))
		// notify scheduler
		//r.NetworkManager.Scheduler.OnQueuedMessage(&m)
	} else {
		r.Log.Printf("Message from %d to %d dropped\n", tm.Sender, tm.Receiver)
	}
}

// check if there is connectivity between two nodes
// returns true if there is connectivity, false otherwise
// if from or to is invalid, it logs an error and returns true
// this is to avoid dropping messages when there is an error in the test
func (r *Router) HasConnectivity(from int, to int) bool {
	if (from < 0 || from >= len(r.RoutingTable)) || (to < 0 || to >= len(r.RoutingTable)) {
		r.Log.Printf("Invalid node IDs: from %d to %d\n", from, to)
		return true
	}
	return r.RoutingTable[from][to]
}

// print routing table in a 2d matrix format
func (r *Router) PrintRoutingTable() {
	for i := 0; i < len(r.RoutingTable); i++ {
		fmt.Printf("RoutingTable[%d]: ", i)
		for j := 0; j < len(r.RoutingTable[i]); j++ {
			fmt.Printf("%t ", r.RoutingTable[i][j])
		}
		fmt.Println()
	}
}

// network faults
// isolate node
func (r *Router) IsolateNode(node int) {
	for i := 0; i < len(r.RoutingTable); i++ {
		r.RoutingTable[node][i] = false
		r.RoutingTable[i][node] = false
	}
	r.PrintRoutingTable()
}

// create network partitions
// accepts a partitions argument, which is a list of sets of node IDs.
// two nodes can communicate with each other if they are on the same partition
// nodes on different partitions cannot communicate with each other
func (r *Router) CreatePartitions(partitions [][]int) {
	// check if partitions are valid - all nodes should be in a single partition
	// and all nodes should be in a partition
	numReplicas := len(r.RoutingTable)
	// print num replicas
	r.Log.Printf("NumReplicas: %d\n", numReplicas)

	visited := make([]int, numReplicas)
	// initialize to -1
	for i := 0; i < numReplicas; i++ {
		visited[i] = -1
	}
	for index, partition := range partitions {
		for _, node := range partition {
			if node < 0 || node >= numReplicas {
				r.Log.Fatalf("Invalid node ID %d\n", node)
			}
			if visited[node] != -1 {
				r.Log.Fatalf("Node %d is in multiple partitions\n", node)
			}
			visited[node] = index
		}
	}

	// create partitions
	for i := 0; i < numReplicas; i++ {
		for j := 0; j < numReplicas; j++ {
			r.RoutingTable[i][j] = visited[i] == visited[j]
		}
	}
	r.PrintRoutingTable()
}
