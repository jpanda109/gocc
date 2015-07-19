package view

import "github.com/nsf/termbox-go"

// NewMessages creates a new Messages object
func NewMessages(x int, y int, w int, h int) *Messages {
	messages := &Messages{
		x,
		y,
		w,
		h,
		make(chan string),
		make([]string, h),
	}
	messages.Start()
	return messages
}

// Messages handles printing messages in its own compartment
type Messages struct {
	x        int
	y        int
	w        int
	h        int
	Incoming chan string
	buffer   []string
}

// Start starts all listeners
func (messages *Messages) Start() {
	go messages.beginListen()
}

// Stop stops listeners and writing etc
func (messages *Messages) Stop() {
	close(messages.Incoming)
}

func (messages *Messages) beginListen() {
	for msg := range messages.Incoming {
		messages.displayMessage(msg)
	}
}

func (messages *Messages) displayMessage(msg string) {
	lines := messages.splitMessage(msg)
	messages.buffer = append(lines, messages.buffer[:messages.h-len(lines)]...)
	curY := messages.y + messages.h - 2
	for _, line := range messages.buffer {
		curX := messages.x
		for i := 0; i < len(line); i++ {
			char := line[i]
			termbox.SetCell(curX, curY, rune(char), termbox.ColorWhite, termbox.ColorBlack)
			curX++
		}
		curY--
	}
	termbox.Flush()
}

func (messages *Messages) splitMessage(msg string) []string {
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
	lines := []string{firstLine}
	for _, v := range otherLines {
		lines = append(lines, "  "+v)
	}
	return lines
}
