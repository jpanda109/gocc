package view

// This file implements the ChatWindow, which simply displays messages
// along with the user's edit bar onto the terminal

import (
	"fmt"
	"log"

	"github.com/jpanda109/gocc/comm"
	"github.com/nsf/termbox-go"
)

// NewChatWindow creates a new chat window
func NewChatWindow() *ChatWindow {
	// Change such that start is public, move width and height, etc into
	// the start function, add listener for window resizing
	window := &ChatWindow{
		0,
		0,
		make(chan []rune),
		make(chan *comm.Message),
		make([]*comm.Message, 0),
	}
	return window
}

// ChatWindow displays current editing buffer and messages
// w is the current width of the terminal
// h is the current height of the terminal
// EditBuffer receives the current editing buffer from a controller
// MsgQ receives messages to display
// msgs is a list of messages needing to be displayed on the terminal
type ChatWindow struct {
	w          int
	h          int
	EditBuffer chan []rune
	MsgQ       chan *comm.Message
	msgs       []*comm.Message
}

// Start begins listeners
func (window *ChatWindow) Start() {
	w, h := termbox.Size()
	window.w, window.h = w, h
	var msgs []*comm.Message
	for i := 0; i < h-2; i++ {
		msgs = append(msgs, &comm.Message{
			Sender: nil,
			Info:   &comm.MsgGob{},
		})
	}
	window.msgs = msgs
	go window.listenEdits()
	go window.listenMsgs()
}

// Stop closes channels and clears screen
func (window *ChatWindow) Stop() {
	close(window.EditBuffer)
	close(window.MsgQ)
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
}

// listenEdits begins listening on the EditBuffer channel, displays onto screen
// accordingly
func (window *ChatWindow) listenEdits() {
	log.Println("Listening to edits")
	for b := range window.EditBuffer {
		for x := 0; x < window.w; x++ {
			termbox.SetCell(x, window.h-1, ' ', termbox.ColorBlack, termbox.ColorBlack)
		}
		for x, ch := range b {
			log.Println(x)
			log.Println(window.h - 1)
			termbox.SetCell(x, window.h-1, ch, termbox.ColorWhite, termbox.ColorBlack)
		}
		termbox.Flush()
	}
}

// listenMsgs begins listening on the MsgQ channel and appends then to the
// msgs buffer accordingly, and then displays the messages onto the terminal
func (window *ChatWindow) listenMsgs() {
	for m := range window.MsgQ {
		window.msgs = append([]*comm.Message{m}, window.msgs[:window.h-3]...)
		y := window.h - 3
		for _, msg := range window.msgs {
			if msg.Sender == nil {
				continue
			}
			for x := 0; x < window.w; x++ {
				termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
			}
			x := 0
			for ; x < len(msg.Sender.Name); x++ {
				termbox.SetCell(x, y, rune(msg.Sender.Name[x]), termbox.ColorCyan, termbox.ColorBlack)
			}
			idPart := fmt.Sprintf(" (%v): ", msg.Sender.ID)
			for i := 0; i < len(idPart); i++ {
				termbox.SetCell(x, y, rune(idPart[i]), termbox.ColorRed, termbox.ColorBlack)
				x++
			}
			body := msg.Info.Body
			for i := 0; i < len(body); i++ {
				termbox.SetCell(x, y, rune(body[i]), termbox.ColorWhite, termbox.ColorBlack)
				x++
			}
			y--
		}
		termbox.Flush()
	}
}
