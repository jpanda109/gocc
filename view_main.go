package main

import (
	"os"

	"github.com/jpanda109/gocc/view"
	"github.com/nsf/termbox-go"
)

func main() {
	termbox.Init()
	defer termbox.Close()
	window := view.NewChatWindow()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	var editBuffer []rune
	for event := range eventQueue {
		if event.Key != 0 {
			key := event.Key
			switch key {
			case termbox.KeyCtrlC:
				termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
				os.Exit(0)
			case termbox.KeyEnter:
				window.MsgQ <- editBuffer
				editBuffer = []rune{}
				window.EditBuffer <- editBuffer
			case termbox.KeyBackspace:
				continue
			}
		} else {
			editBuffer = append(editBuffer, event.Ch)
			window.EditBuffer <- editBuffer
		}
	}
}

// func main() {
// 	termbox.Init()
// 	defer termbox.Close()
// 	w, h := termbox.Size()
// 	input := view.NewChatInput(0, 0, w, h)
// 	messages := view.NewMessages(0, 1, w, h)
// 	eventQueue := make(chan termbox.Event)
// 	go func() {
// 		for {
// 			eventQueue <- termbox.PollEvent()
// 		}
// 	}()
// 	go func() {
// 		for event := range eventQueue {
// 			if event.Key != 0 {
// 				key := event.Key
// 				switch key {
// 				case termbox.KeyCtrlC:
// 					input.Stop()
// 					messages.Stop()
// 					termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
// 				default:
// 					input.IncomingKey <- key
// 				}
// 			} else {
// 				input.IncomingCh <- event.Ch
// 			}
// 		}
// 	}()
// 	for msg := range input.OutgoingMessages {
// 		messages.Incoming <- msg
// 	}
// }
