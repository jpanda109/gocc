package view

import "github.com/nsf/termbox-go"

// NewInput creates a new Input struct
func NewInput(x int, y int, w int, h int) *Input {
	input := &Input{
		x,
		y,
		w,
		h,
		make(chan rune),
		make(chan termbox.Key),
		make(chan string),
		make([]rune, 0),
		make(chan bool),
	}
	return input
}

// Input displays input to screen
type Input struct {
	x                int
	y                int
	w                int
	h                int
	IncomingCh       chan rune
	IncomingKey      chan termbox.Key
	OutgoingMessages chan string
	buffer           []rune
	quit             chan bool
}

// Start begins input handling and displaying
func (input *Input) Start() {
	go input.handleInput()
	<-input.quit
}

// Stop stops the input from doing anything
func (input *Input) Stop() {
	input.quit <- true
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.Flush()
	close(input.OutgoingMessages)
}

func (input *Input) handleInput() {
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

func (input *Input) displayBuffer() {
	for i := 0; i < input.w; i++ {
		termbox.SetCell(input.x+i, input.y, ' ', termbox.ColorWhite, termbox.ColorBlack)
	}
	for i, ch := range input.buffer {
		termbox.SetCell(input.x+i, input.y, ch, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}
