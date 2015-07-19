package comm

import (
	"net"
	"strings"
)

// NewChatRoom creates and returns a chat room reference
func NewChatRoom() *ChatRoom {
	chatRoom := &ChatRoom{
		make([]*Client, 0),
		[]string{"127.0.0.1"},
		make(chan net.Conn),
		make(chan string),
		make(chan string),
		make(chan bool),
	}
	chatRoom.Start()
	return chatRoom
}

// ChatRoom handles broadcasting messages to group of containers
type ChatRoom struct {
	clients        []*Client
	Whitelist      []string
	NewConnections chan net.Conn
	Incoming       chan string
	Outgoing       chan string
	Quit           chan bool
}

// Start sets listener, reader, and writer
func (cr *ChatRoom) Start() {
	go cr.beginListen()
	go cr.beginWrite()
}

func (cr *ChatRoom) beginListen() {
	for c := range cr.NewConnections {
		client := NewClient(c)
		if cr.whitelisted(c.RemoteAddr().String()) {
			go func() {
				for {
					cr.Incoming <- <-client.Incoming
				}
			}()
			cr.clients = append(cr.clients, client)
			// cr.Broadcast("New Client: " + client.Name)
		} else {
			c.Close()
		}
	}
}

func (cr *ChatRoom) beginWrite() {
	for msg := range cr.Outgoing {
		cr.Broadcast(msg)
	}
}

// Broadcast sends message to all clients
func (cr *ChatRoom) Broadcast(msg string) {
	for _, client := range cr.clients {
		client.Outgoing <- msg
	}
}

func (cr *ChatRoom) whitelisted(addr string) bool {
	for _, a := range cr.Whitelist {
		if strings.HasPrefix(addr, a) {
			return true
		}
	}
	return false
}
