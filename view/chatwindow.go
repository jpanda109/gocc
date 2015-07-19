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
		NewMessages(0, 0, w-1, h-1),
		NewChatInput(x, y+h-2, w, 2),
		make(chan string),
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
	IncomingMsgs chan string
	OutgoingMsgs chan string
	Events       chan termbox.Event
	quit         chan bool
}

// Start starts the chat window processing
func (window *ChatWindow) Start() {
	go window.listenEvents()
	go window.listenMsgs()
	go window.listenInput()
}

// Stop stops the chat window processing
func (window *ChatWindow) Stop() {
	window.messages.Stop()
	window.chatInput.Stop()
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
	close(window.IncomingMsgs)
	close(window.OutgoingMsgs)
	close(window.Events)
	close(window.quit)
}

func (window *ChatWindow) listenEvents() {
	for event := range window.Events {
		switch {
		case event.Key != 0:
			window.chatInput.IncomingKey <- event.Key
		case event.Ch != 0:
			window.chatInput.IncomingCh <- event.Ch
		}
	}
}

func (window *ChatWindow) listenInput() {
	for msg := range window.chatInput.OutgoingMessages {
		window.OutgoingMsgs <- msg
	}
}

func (window *ChatWindow) listenMsgs() {
	for {
		window.messages.Incoming <- <-window.IncomingMsgs
	}
}

// Quit returns control to terminal
func (window *ChatWindow) Quit() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
}
