package input

import (
	"sync"

	"github.com/jpanda109/gocc/comm"
	"github.com/jpanda109/gocc/view"
	"github.com/nsf/termbox-go"
)

// NewHandler returns reference to Handler object
func NewHandler(addr, name string) *Handler {
	chatroom := comm.NewChatRoom()
	handler := &Handler{
		comm.NewConnHandler(addr, name, chatroom),
		chatroom,
		view.NewChatWindow(),
		make([]rune, 0),
	}
	return handler
}

// Handler handles input and controls commucation between peers
// and views
type Handler struct {
	cHandler   *comm.ConnHandler
	chatroom   *comm.ChatRoom
	window     *view.ChatWindow
	editBuffer []rune
}

// Start begins handler listening for input
func (h *Handler) Start() *sync.WaitGroup {
	var wg sync.WaitGroup
	go h.handleConns()
	go h.handleMessages()
	wg.Add(1)
	go h.handleEvents(&wg)
	return &wg
}

// Connect connects to chat room and adds all existing peers
func (h *Handler) Connect(addr string) {
	peers := h.cHandler.Dial(addr)
	for _, peer := range peers {
		h.chatroom.AddPeer(peer)
	}
}

func (h *Handler) handleMessages() {
	for {
		msg := h.chatroom.Receive()
		h.window.MsgQ <- []rune(msg.Body)
	}
}

func (h *Handler) handleConns() {
	h.cHandler.Listen()
	for {
		peer := h.cHandler.GetPeer()
		h.chatroom.AddPeer(peer)
	}
}

func (h *Handler) handleEvents(wg *sync.WaitGroup) {
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
				h.window.Stop()
				return
			case termbox.KeyEnter:
				h.chatroom.Broadcast(string(h.editBuffer))
				h.window.MsgQ <- []rune(h.editBuffer)
				h.editBuffer = []rune{}
				h.window.EditBuffer <- h.editBuffer
			case termbox.KeyBackspace:
				if curlen := len(h.editBuffer); curlen > 0 {
					h.editBuffer = append([]rune{}, h.editBuffer[:curlen-1]...)
					h.window.EditBuffer <- h.editBuffer
				}
			case termbox.KeySpace:
				h.editBuffer = append(h.editBuffer, ' ')
				h.window.EditBuffer <- h.editBuffer
			}
		} else {
			h.editBuffer = append(h.editBuffer, event.Ch)
			h.window.EditBuffer <- h.editBuffer
		}
	}
}
