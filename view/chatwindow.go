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
		NewMessages(0, 0, x-1, y-1),
		make(chan string),
		make(chan bool),
	}
	chatWindow.Start()
	return chatWindow
}

// ChatWindow handles displaying interface needed for chat
type ChatWindow struct {
	x            int
	y            int
	messages     *Messages
	incomingMsgs chan string
	quit         chan bool
}

// Start starts the chat window processing
func (window *ChatWindow) Start() {
	window.messages.Start()
	go window.listenMsgs()
	<-window.quit
}

// Stop stops the chat window processing
func (window *ChatWindow) Stop() {
	window.quit <- true
	window.messages.Stop()
}

func (window *ChatWindow) listenMsgs() {
	for {
		window.messages.Incoming <- <-window.incomingMsgs
	}
}

// Quit returns control to terminal
func (window *ChatWindow) Quit() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
	termbox.Close()
}
