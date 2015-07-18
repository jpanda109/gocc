package main

import (
	"fmt"

	"github.com/jpanda109/gocc/view"
	"github.com/nsf/termbox-go"
)

func main() {
	termbox.Init()
	w, h := termbox.Size()
	input := view.NewInput(0, 0, w, h)
	go input.Start()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	go func() {
		for event := range eventQueue {
			if event.Key != 0 {
				if event.Key == termbox.KeyCtrlC {
					input.Stop()
					break
				} else {
					input.IncomingKey <- event.Key
				}
			} else {
				input.IncomingCh <- event.Ch
			}
		}
	}()
	for msg := range input.OutgoingMessages {
		fmt.Println(msg)
	}
	termbox.Close()
}
