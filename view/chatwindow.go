package view

import (
	"github.com/nsf/termbox-go"
)

// NewChatWindow creates a new chat window
func NewChatWindow() *ChatWindow {
	w, h := termbox.Size()
	window := &ChatWindow{
		w,
		h,
		make(chan []rune),
		make(chan []rune),
		make([][]rune, h-2),
	}
	window.start()
	return window
}

// ChatWindow displays current editing buffer and messages
type ChatWindow struct {
	w          int
	h          int
	EditBuffer chan []rune
	MsgQ       chan []rune
	msgs       [][]rune
}

func (window *ChatWindow) start() {
	go window.listenEdits()
	go window.listenMsgs()
}

func (window *ChatWindow) listenEdits() {
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

func (window *ChatWindow) listenMsgs() {
	for m := range window.MsgQ {
		window.msgs = append([][]rune{m}, window.msgs[:window.h-3]...)
		y := window.h - 3
		for _, line := range window.msgs {
			for x := 0; x < window.w; x++ {
				termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
			}
			for x := 0; x < len(line); x++ {
				termbox.SetCell(x, y, rune(line[x]), termbox.ColorWhite, termbox.ColorBlack)
			}
			y--
		}
		termbox.Flush()
	}
}
