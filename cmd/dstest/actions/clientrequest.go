package actions

type ClientRequestAction struct {
	Name string
}

// make sure ClientRequestAction implements the Action interface
var _ Action = (*ClientRequestAction)(nil)

func (cra *ClientRequestAction) GetType() ActionType {
	return ClientRequest
}
