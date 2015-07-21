package notmain2

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
				window.Stop()
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
