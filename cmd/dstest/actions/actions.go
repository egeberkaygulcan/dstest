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

func NewAction(actionType ActionType, params map[string]interface{}) Action {
	switch actionType {
	case SendMessage:
		return &DeliverMessageAction{
			Sender:   params["sender"].(int),
			Receiver: params["receiver"].(int),
			Name:     params["name"].(string),
		}
	case InjectFault:
		return &InjectFaultAction{
			// TODO: Implement
		}
	case ClientRequest:
		return &ClientRequestAction{
			ClientId: params["clientId"].(int),
		}
	default:
		return nil
	}
}
