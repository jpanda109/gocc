package input

// This file defines the Controller, which handles user input and
// actions received from peers

import (
	"log"
	"strings"
	"sync"

	"github.com/jpanda109/gocc/comm"
	"github.com/jpanda109/gocc/config"
	"github.com/jpanda109/gocc/view"
	"github.com/nsf/termbox-go"
)

type ownPeer struct {
	addr string
	name string
}

func (p *ownPeer) ID() int {
	return 0
}

func (p *ownPeer) Name() string {
	return p.name
}

func (p *ownPeer) Send(msg *comm.MsgGob) error {
	return nil
}

func (p *ownPeer) Receive() (*comm.MsgGob, error) {
	return nil, nil
}

func (p *ownPeer) String() string {
	return p.addr + "," + p.name
}

// NewController returns reference to Controller object
func NewController(addr, name string) *Controller {
	chatroom := comm.NewChatRoom()
	controller := &Controller{
		&ownPeer{addr, name},
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
// self is a reference to a peer representing the user
// editBuffer is the current buffer that the user sends upon choosing so
// quit is a channel signalling when the user decides to stop the application
type Controller struct {
	self       comm.Peer
	cHandler   *comm.ConnHandler
	chatroom   comm.ChatRoom
	window     *view.ChatWindow
	editBuffer []rune
	quit       chan bool
}

// Start begins handler listening for input
// returns a WaitGroup which tells anyone using the controller when the user
// 	decides to stop running the application
func (c *Controller) Start() *sync.WaitGroup {
	log.Println("Controller.Start() called")
	c.window.Start()
	var wg sync.WaitGroup
	go c.handleConns()
	go c.handleMessages()
	wg.Add(1)
	go c.listenEvents(&wg)
	return &wg
}

// Connect connects to chat room and adds all existing peers
// returns an error if there's an issue connecting to the peer
func (c *Controller) Connect(addr string) error {
	peers, err := c.cHandler.Dial(addr)
	if err != nil {
		return err
	}
	for _, peer := range peers {
		c.chatroom.AddPeer(peer)
	}
	return nil
}

// handleMessages dictates how other peers' messages are treated
func (c *Controller) handleMessages() {
	for {
		msg := c.chatroom.Receive()
		if msg.Info.Action == comm.Public {
			c.window.MsgQ <- &view.Message{
				SenderID:   msg.SenderID,
				SenderName: msg.SenderName,
				Type:       view.Public,
				Body:       msg.Info.Body,
			}
		}
	}
}

// handleConns dictates how new connections are treated
func (c *Controller) handleConns() {
	c.cHandler.Listen()
	for {
		peer := c.cHandler.Peer()
		c.chatroom.AddPeer(peer)
	}
}

// listenEvents gathers input from termbox.PollEvent()
func (c *Controller) listenEvents(wg *sync.WaitGroup) {
	defer wg.Done()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	go c.handleEvents(eventQueue)
	log.Println("controller waiting on quit")
	<-c.quit
	log.Println("controller no longer listening")
}

func (c *Controller) handleInput() {
	if len(c.editBuffer) == 0 {
		return
	}
	tokens := strings.Split(string(c.editBuffer), " ")
	if tokens[0][0] == '/' {
		switch tokens[0] {
		case "/save":
		case "/whitelistfriends":
			for _, f := range config.Friends() {
				c.cHandler.Whitelist(f.Addr)
			}
		case "/whitelist":
			addrs := tokens[1:]
			for _, addr := range addrs {
				c.cHandler.Whitelist(addr)
			}
		case "/friends":
		}
	} else {
		c.chatroom.Broadcast(string(c.editBuffer))
		c.window.MsgQ <- &view.Message{
			SenderID:   c.self.ID(),
			SenderName: c.self.Name(),
			Type:       view.Public,
			Body:       string(c.editBuffer),
		}
	}
	c.editBuffer = []rune{}
	c.window.EditBuffer <- c.editBuffer
}

// handleEvents dictates how user input is treated and what actions to take
func (c *Controller) handleEvents(eventQueue chan termbox.Event) {
	for event := range eventQueue {
		if event.Key != 0 {
			switch event.Key {
			case termbox.KeyCtrlC:
				c.window.Stop()
				c.quit <- true
			case termbox.KeyEnter:
				c.handleInput()
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
