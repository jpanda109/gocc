package view

import (
	"github.com/nsf/termbox-go"
)

// ChatWindow handles displaying interface needed for chat
type ChatWindow struct {
	x int
	y int
}

// Init initializes the chat window
func (window *ChatWindow) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	window.x, window.y = termbox.Size()
	return nil
}

// Quit returns control to terminal
func (window *ChatWindow) Quit() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
	termbox.Close()
}
