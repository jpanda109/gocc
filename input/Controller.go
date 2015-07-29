package input

import (
	"sync"

	"github.com/jpanda109/gocc/comm"
	"github.com/jpanda109/gocc/view"
	"github.com/nsf/termbox-go"
)

// NewController returns reference to Handler object
func NewController(addr, name string) *Controller {
	chatroom := comm.NewChatRoom()
	controller := &Controller{
		comm.NewConnHandler(addr, name, chatroom),
		chatroom,
		view.NewChatWindow(),
		make([]rune, 0),
	}
	return controller
}

// Controller handles input and controls commucation between peers
// and views
type Controller struct {
	cHandler   *comm.ConnHandler
	chatroom   *comm.ChatRoom
	window     *view.ChatWindow
	editBuffer []rune
}

// Start begins handler listening for input
func (c *Controller) Start() *sync.WaitGroup {
	var wg sync.WaitGroup
	go c.handleConns()
	go c.handleMessages()
	wg.Add(1)
	go c.handleEvents(&wg)
	return &wg
}

// Connect connects to chat room and adds all existing peers
func (c *Controller) Connect(addr string) {
	peers := c.cHandler.Dial(addr)
	for _, peer := range peers {
		c.chatroom.AddPeer(peer)
	}
}

func (c *Controller) handleMessages() {
	for {
		msg := c.chatroom.Receive()
		c.window.MsgQ <- []rune(msg.Body)
	}
}

func (c *Controller) handleConns() {
	c.cHandler.Listen()
	for {
		peer := c.cHandler.GetPeer()
		c.chatroom.AddPeer(peer)
	}
}

func (c *Controller) handleEvents(wg *sync.WaitGroup) {
	defer wg.Done()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	for event := range eventQueue {
		if event.Key != 0 {
			switch event.Key {
			case termbox.KeyCtrlC:
				c.window.Stop()
				return
			case termbox.KeyEnter:
				c.chatroom.Broadcast(string(c.editBuffer))
				c.window.MsgQ <- []rune(c.editBuffer)
				c.editBuffer = []rune{}
				c.window.EditBuffer <- c.editBuffer
			case termbox.KeyBackspace:
				if curlen := len(c.editBuffer); curlen > 0 {
					c.editBuffer = append([]rune{}, c.editBuffer[:curlen-1]...)
					c.window.EditBuffer <- c.editBuffer
				}
			case termbox.KeySpace:
				c.editBuffer = append(c.editBuffer, ' ')
				c.window.EditBuffer <- c.editBuffer
			}
		} else {
			c.editBuffer = append(c.editBuffer, event.Ch)
			c.window.EditBuffer <- c.editBuffer
		}
	}
}
