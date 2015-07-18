package view

import (
	"github.com/nsf/termbox-go"
)

// NewChatWindow creates a new chat window
func NewChatWindow(x int, y int, w int, h int) *ChatWindow {
	chatWindow := &ChatWindow{
		x,
		y,
		w,
		h,
		NewMessages(0, 0, x-1, y-1),
		NewChatInput(x, y+h-2, w, 2),
		make(chan string),
		make(chan termbox.Event),
		make(chan bool),
	}
	chatWindow.Start()
	return chatWindow
}

// ChatWindow handles displaying and gathering use input interface needed for chat
type ChatWindow struct {
	x            int
	y            int
	w            int
	h            int
	messages     *Messages
	chatInput    *ChatInput
	incomingMsgs chan string
	events       chan termbox.Event
	quit         chan bool
}

// Start starts the chat window processing
func (window *ChatWindow) Start() {
	window.messages.Start()
	go window.listenInput()
	go window.listenMsgs()
	<-window.quit
}

// Stop stops the chat window processing
func (window *ChatWindow) Stop() {
	window.quit <- true
	window.messages.Stop()
	window.chatInput.Stop()
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
}

func (window *ChatWindow) listenInput() {
	for {
		window.events <- termbox.PollEvent()
	}
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
