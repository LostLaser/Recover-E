package communication

type message interface {
	sealedMessage()
}

// Event contains information for an action regarding two servers
type Event struct {
	From        string `json:"from"`
	To          string `json:"to"`
	ActionType  Action `json:"action"`
	MessageType string `json:"type"`
}

// Control contains information for an action concerning a single server
type Control struct {
	Target      string `json:"target"`
	ActionType  Action `json:"action"`
	MessageType string `json:"type"`
}

// NewEvent creates a new Event message
func NewEvent(from, to string, action Action) Event {
	e := new(Event)
	e.From = from
	e.To = to
	e.ActionType = action
	e.MessageType = "EVENT"

	return *e
}

func (e Event) sealedMessage() {}

// NewControl creates a new Control message
func NewControl(target string, action Action) Control {
	c := new(Control)
	c.Target = target
	c.ActionType = action
	c.MessageType = "CONTROL"

	return *c
}

func (c Control) sealedMessage() {}
