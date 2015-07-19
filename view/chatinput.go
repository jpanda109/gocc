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
		make(chan termbox.Key),
		make(chan string),
		make([]rune, 0),
	}
	input.Start()
	return input
}

// ChatInput displays input to screen
type ChatInput struct {
	x                int
	y                int
	w                int
	h                int
	IncomingCh       chan rune
	IncomingKey      chan termbox.Key
	OutgoingMessages chan string
	buffer           []rune
}

// Start begins input handling and displaying
func (input *ChatInput) Start() {
	go input.handleInput()
}

// Stop stops the input from doing anything
func (input *ChatInput) Stop() {
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
		case key := <-input.IncomingKey:
			switch key {
			case termbox.KeyEnter:
				if len(input.buffer) < input.w {
					input.OutgoingMessages <- string(input.buffer[:])
					input.buffer = make([]rune, 0)
				}
			case termbox.KeyBackspace:
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
