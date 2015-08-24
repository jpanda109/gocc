package view

// This file implements the ChatWindow, which simply displays messages
// along with the user's edit bar onto the terminal

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

// MType is an alias which defines different types of messages which can be
// displayed
type MType int

const (
	// Public denotes messages which are broadcasted to all peers
	Public MType = iota
	// Internal denotes messages which are only viewable by the person who
	// sent them
	Internal
)

// Message defines a struct which contains the necessary information to display
// text correctly
type Message struct {
	SenderID   int
	SenderName string
	Type       MType
	Body       string
}

// NewChatWindow creates a new chat window
func NewChatWindow() *ChatWindow {
	// Change such that start is public, move width and height, etc into
	// the start function, add listener for window resizing
	window := &ChatWindow{
		0,
		0,
		make(chan []rune),
		make(chan *Message),
		make([]*Message, 0),
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
	MsgQ       chan *Message
	msgs       []*Message
}

// Start begins listeners
func (window *ChatWindow) Start() {
	w, h := termbox.Size()
	window.w, window.h = w, h
	var msgs []*Message
	for i := 0; i < h-2; i++ {
		msgs = append(msgs, &Message{
			SenderID:   -1,
			SenderName: "",
			Type:       Public,
			Body:       "",
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
			termbox.SetCell(x, window.h-1, ch, termbox.ColorWhite, termbox.ColorBlack)
		}
		termbox.Flush()
	}
}

// listenMsgs begins listening on the MsgQ channel and appends then to the
// msgs buffer accordingly, and then displays the messages onto the terminal
func (window *ChatWindow) listenMsgs() {
	for m := range window.MsgQ {
		window.msgs = append([]*Message{m}, window.msgs[:window.h-3]...)
		y := window.h - 3
		for _, msg := range window.msgs {
			if msg.SenderID == -1 {
				continue
			}
			if msg.Type == Public {
				for x := 0; x < window.w; x++ {
					termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
				}
				x := 0
				for ; x < len(msg.SenderName); x++ {
					termbox.SetCell(x, y, rune(msg.SenderName[x]), termbox.ColorCyan, termbox.ColorBlack)
				}
				idPart := fmt.Sprintf(" (%v): ", msg.SenderID)
				for i := 0; i < len(idPart); i++ {
					termbox.SetCell(x, y, rune(idPart[i]), termbox.ColorRed, termbox.ColorBlack)
					x++
				}
				body := msg.Body
				for i := 0; i < len(body); i++ {
					termbox.SetCell(x, y, rune(body[i]), termbox.ColorWhite, termbox.ColorBlack)
					x++
				}
				y--
			} else if msg.Type == Internal {
				for x := 0; x < window.w; x++ {
					termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
				}
				for x := 0; x < len(msg.Body); x++ {
					termbox.SetCell(x, y, rune(msg.Body[x]), termbox.ColorYellow, termbox.ColorBlack)
				}
			}
		}
		termbox.Flush()
	}
}
