package comm

import "fmt"

// Message contains information about a message received from a peer
// Sender is a pointer to the local peer instance which sent it
// info contains the message infomation such as the action type and body
type Message struct {
	SenderID   int
	SenderName string
	Info       *MsgGob
}

// String returns a human readable form of a Message
func (msg *Message) String() string {
	return fmt.Sprintf("%s (%s) > %s", msg.SenderName, msg.Info.Action,
		msg.Info.Body)
}
