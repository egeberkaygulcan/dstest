package actions

type ActionType string

const (
	SendMessage   ActionType = "SendMessage"
	InjectFault   ActionType = "InjectFault"
	ClientRequest ActionType = "ClientRequest"
)

func (at ActionType) String() string {
	switch at {
	case SendMessage:
		return "SendMessage"
	case InjectFault:
		return "InjectFault"
	case ClientRequest:
		return "ClientRequest"
	default:
		return "Unknown"
	}
}

type Action interface {
	GetType() ActionType
}

func NewAction(actionType ActionType, params map[string]string) Action {
	switch actionType {
	case SendMessage:
		return &DeliverMessageAction{}
	case InjectFault:
		return &InjectFaultAction{}
	case ClientRequest:
		return &ClientRequestAction{}
	default:
		return nil
	}
}
