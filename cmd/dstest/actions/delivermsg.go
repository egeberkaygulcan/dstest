package actions

type DeliverMessageAction struct {
	Sender   int
	Receiver int
	Name     string
}

// make sure DeliverMessageAction implements the Action interface
var _ Action = (*DeliverMessageAction)(nil)

func (dma *DeliverMessageAction) GetType() ActionType {
	return SendMessage
}
