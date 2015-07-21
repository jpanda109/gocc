package comm

// NewChatRoom creates and returns a chat room reference
func NewChatRoom() *ChatRoom {
	chatRoom := &ChatRoom{
		make([]*Client, 0),
		make(chan *Client),
		make(chan string),
		make(chan string),
	}
	chatRoom.Start()
	return chatRoom
}

// ChatRoom handles broadcasting messages to group of containers
type ChatRoom struct {
	clients    []*Client
	NewClients chan *Client
	Incoming   chan string
	Outgoing   chan string
}

// Start sets listener, reader, and writer
func (cr *ChatRoom) Start() {
	go cr.beginListen()
	go cr.beginWrite()
}

func (cr *ChatRoom) beginListen() {
	for c := range cr.NewClients {
		cr.clients = append(cr.clients, c)
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
