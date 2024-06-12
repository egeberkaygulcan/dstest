package scheduling

//import "github.com/egeberkaygulcan/dstest/cmd/dstest/network"

type Scheduler interface {
	Init()
	OnQueuedMessage(m *any)
	OnStartup()
	OnShutdown()
}
