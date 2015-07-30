package view

import (
	"fmt"
	"log"

	"github.com/jpanda109/gocc/comm"
	"github.com/nsf/termbox-go"
)

// NewChatWindow creates a new chat window
func NewChatWindow() *ChatWindow {
	w, h := termbox.Size()
	var msgs []*comm.Message
	for i := 0; i < h-2; i++ {
		msgs = append(msgs, &comm.Message{
			Sender: nil,
			Body:   "",
		})
	}
	window := &ChatWindow{
		w,
		h,
		make(chan []rune),
		make(chan *comm.Message),
		msgs,
	}
	window.start()
	return window
}

// ChatWindow displays current editing buffer and messages
type ChatWindow struct {
	w          int
	h          int
	EditBuffer chan []rune
	MsgQ       chan *comm.Message
	msgs       []*comm.Message
}

func (window *ChatWindow) start() {
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
				termbox.SetCell(x, y, rune(idPart[i]), termbox.ColorBlue, termbox.ColorBlack)
				x++
			}
			for i := 0; i < len(msg.Body); i++ {
				termbox.SetCell(x, y, rune(msg.Body[i]), termbox.ColorWhite, termbox.ColorBlack)
				x++
			}
			y--
		}
		termbox.Flush()
	}
}
