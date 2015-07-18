package view

import "github.com/nsf/termbox-go"

// NewChatInput creates a new Input struct
func NewChatInput(x int, y int, w int, h int) *ChatInput {
	input := &ChatInput{
		x,
		y,
		w,
		h,
		make(chan rune),
		make(chan Command),
		make(chan string),
		make([]rune, 0),
		make(chan bool),
	}
	return input
}

// Command is command to be sent to chat input
type Command int

const (
	// Backspace deletes last char
	Backspace Command = iota
	// Submit deletes and sends buffer to outgoing messages channel
	Submit
)

// ChatInput displays input to screen
type ChatInput struct {
	x                int
	y                int
	w                int
	h                int
	IncomingCh       chan rune
	IncCommand       chan Command
	OutgoingMessages chan string
	buffer           []rune
	quit             chan bool
}

// Start begins input handling and displaying
func (input *ChatInput) Start() {
	go input.handleInput()
	<-input.quit
}

// Stop stops the input from doing anything
func (input *ChatInput) Stop() {
	input.quit <- true
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.Flush()
	close(input.OutgoingMessages)
}

func (input *ChatInput) handleInput() {
	for {
		select {
		case ch := <-input.IncomingCh:
			if len(input.buffer) < input.w {
				input.buffer = append(input.buffer, ch)
			}
		case command := <-input.IncCommand:
			switch command {
			case Submit:
				if len(input.buffer) < input.w {
					input.OutgoingMessages <- string(input.buffer[:])
					input.buffer = make([]rune, 0)
				}
			case Backspace:
				if len(input.buffer) > 0 {
					input.buffer = input.buffer[:len(input.buffer)-1]
				}
			}
		}
		input.displayBuffer()
	}
}

func (input *ChatInput) displayBuffer() {
	for i := 0; i < input.w; i++ {
		termbox.SetCell(input.x+i, input.y, ' ', termbox.ColorWhite, termbox.ColorBlack)
	}
	for i, ch := range input.buffer {
		termbox.SetCell(input.x+i, input.y, ch, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}
