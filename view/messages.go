package view

import (
	"github.com/nsf/termbox-go"
)

// NewMessages creates a new Messages object
func NewMessages(x int, y int, w int, h int) *Messages {
	messages := &Messages{
		x,
		y,
		w,
		h,
		make(chan string),
		make([]string, 0),
	}
	return messages
}

// Messages handles printing messages in its own compartment
type Messages struct {
	x        int
	y        int
	w        int
	h        int
	incoming chan string
	buffer   []string
}

func (messages *Messages) displayMessage(msg string) {
	firstLine := ""
	remainderMsg := ""
	if len(msg) > messages.w {
		firstLine = msg[0:messages.w]
		remainderMsg = msg[messages.w:len(msg)]
	} else {
		firstLine = msg
	}
	var otherLines []string
	for i := 0; i < len(remainderMsg)/(messages.w-4); i++ {
		otherLines = append(otherLines, remainderMsg[i*(messages.w-4):(i+1)*(messages.w-4)])
	}
	if remLen := len(remainderMsg) % (messages.w - 4); remLen > 0 {
		otherLines = append(otherLines, remainderMsg[len(remainderMsg)-remLen-1:len(remainderMsg)])
	}
	lines := (len(msg)-messages.w)/(messages.w-4) + 1
	termbox.Flush()
}
