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
		comm.NewPeer(nil, addr, name),
		comm.NewConnHandler(addr, name, chatroom),
		chatroom,
		view.NewChatWindow(),
		make([]rune, 0),
		make(chan bool),
	}
	return controller
}

// Controller handles input and controls commucation between peers
// and views
type Controller struct {
	self       *comm.Peer
	cHandler   *comm.ConnHandler
	chatroom   *comm.ChatRoom
	window     *view.ChatWindow
	editBuffer []rune
	quit       chan bool
}

// Start begins handler listening for input
func (c *Controller) Start() *sync.WaitGroup {
	var wg sync.WaitGroup
	go c.handleConns()
	go c.handleMessages()
	wg.Add(1)
	go c.listenEvents(&wg)
	return &wg
}

// Connect connects to chat room and adds all existing peers
func (c *Controller) Connect(addr string) {
	peers, err := c.cHandler.Dial(addr)
	if err != nil {
		c.window.Stop()
		c.quit <- true
	}
	for _, peer := range peers {
		c.chatroom.AddPeer(peer)
	}
}

func (c *Controller) handleMessages() {
	for {
		msg := c.chatroom.Receive()
		c.window.MsgQ <- msg
	}
}

func (c *Controller) handleConns() {
	c.cHandler.Listen()
	for {
		peer := c.cHandler.GetPeer()
		c.chatroom.AddPeer(peer)
	}
}

func (c *Controller) listenEvents(wg *sync.WaitGroup) {
	defer wg.Done()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	go c.handleEvents(eventQueue)
	<-c.quit
}

func (c *Controller) handleEvents(eventQueue chan termbox.Event) {
	for event := range eventQueue {
		if event.Key != 0 {
			switch event.Key {
			case termbox.KeyCtrlC:
				c.window.Stop()
				c.quit <- true
			case termbox.KeyEnter:
				c.chatroom.Broadcast(string(c.editBuffer))
				c.window.MsgQ <- &comm.Message{
					Sender: c.self,
					Body:   string(c.editBuffer),
				}
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
