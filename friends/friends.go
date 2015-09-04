package friends

// The friends package implements a feature of the gocc app which will
// display and manipulate saved friends' info along with being able to
// transition to the main application by connecting to a selected friend.

import (
	"github.com/jpanda109/gocc/app"
	"github.com/jpanda109/gocc/config"
	"github.com/nsf/termbox-go"
)

// FriendApp contains the necessary information for the friends feature
type FriendApp struct {
	manager  *app.Manager
	listView *ListView

	quit chan bool
}

// Start starts this feature
func (app *FriendApp) Start() {
	app.quit = make(chan bool)
	app.listView = NewListView()
	w, h := termbox.Size()
	app.listView.Start(w, h, config.Friends())
	go app.listenEvents()
}

// Stop tells the current app feature to stop to give control back to manager
func (app *FriendApp) Stop() {
	app.quit <- true
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
}

// SetManager is used to follow interface app.Screen
func (app *FriendApp) SetManager(manager *app.Manager) {
	app.manager = manager
}

func (app *FriendApp) listenEvents() {
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	go app.handleEvents(eventQueue)
	<-app.quit
}

func (app *FriendApp) handleEvents(events chan termbox.Event) {
	for event := range events {
		if event.Key != 0 {
			switch event.Key {
			case termbox.KeyCtrlC:
				app.manager.Quit()
				break
			}
		}
	}
}
