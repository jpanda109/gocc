package comm

import (
	"net"
)

// NewChatRoom creates and returns a chat room pointer with given password
func NewChatRoom(password string) *ChatRoom {
	chatRoom := &ChatRoom{
		password,
		make([]*Client, 0),
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
	password       string
	clients        []*Client
	NewConnections chan net.Conn
	Incoming       chan string
	Outgoing       chan string
	Quit           chan bool
}

// Start sets listener, reader, and writer
func (cr *ChatRoom) Start() {
	go cr.beginListen()
	go cr.beginWrite()
	<-cr.Quit
}

func (cr *ChatRoom) beginListen() {
	for c := range cr.NewConnections {
		client := NewClient(c)
		clientPass := <-client.Incoming
		if clientPass == cr.password {
			client.Name = <-client.Incoming
			go func() {
				for {
					cr.Incoming <- <-client.Incoming
				}
			}()
			cr.clients = append(cr.clients, client)
			cr.Broadcast("New Client: " + client.Name)
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
