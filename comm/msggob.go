package comm

// MsgGob is a gob which used to send necessary data through tcp sockets
// action is the type of message being sent
// body is the body of the message being sent
type MsgGob struct {
	Action Action
	Body   string
}

// Action represents the type of action that the message represents
type Action int

const (
	// Public indicates the message is broadcasted to all
	Public Action = iota
	// Private indicates the message was sent to one peer
	Private
	// Whitelist indicates a modification to the chatroom whitelist
	Whitelist
)
