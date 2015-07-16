package view

import (
	"github.com/nsf/termbox-go"
)

// NewChatWindow creates a new chat window
func NewChatWindow() *ChatWindow {
	termbox.Init()
	x, y := termbox.Size()
	chatWindow := &ChatWindow{
		x,
		y,
		NewMessages(0, 0, x, y),
		make(chan string),
	}
	return chatWindow
}

// ChatWindow handles displaying interface needed for chat
type ChatWindow struct {
	x            int
	y            int
	messages     *Messages
	incomingMsgs chan string
}

// Start starts the chat window processing
func (window *ChatWindow) Start() {
	window.messages.Start()
}

func (window *ChatWindow) listenMsgs() {

}

// Quit returns control to terminal
func (window *ChatWindow) Quit() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
	termbox.Close()
}
