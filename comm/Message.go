package comm

import "fmt"

// Message contains information about a message received from a peer
type Message struct {
	Sender *Peer
	Body   string
}

// String returns a human readable form of a Message
func (msg *Message) String() string {
	return fmt.Sprintf("%s > %s", msg.Sender, msg.Body)
}
