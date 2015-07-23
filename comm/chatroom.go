package comm

import "net"

// NewChatRoom creates and returns a pointer to chat room
func NewChatRoom() *ChatRoom {
	room := &ChatRoom{
		make(chan peerRemoval),
		make(chan peerAddition),
		make(chan string),
		make(chan string),
		[]*Peer{},
	}
	room.start()
	return room
}

// ChatRoom sends and receives messages to peers
type ChatRoom struct {
	removalQueue  chan peerRemoval
	additionQueue chan peerAddition
	Incoming      chan string
	Outgoing      chan string
	peers         []*Peer
}

// AddPeer adds peer to chat room with given info
func (room *ChatRoom) AddPeer(conn net.Conn, addr string, name string) {
	room.additionQueue <- peerAddition{conn, addr, name}
}

// RemovePeer removes peer from chat room with given info
func (room *ChatRoom) RemovePeer(addr string, name string) {
	room.removalQueue <- peerRemoval{addr, name}
}

func (room *ChatRoom) start() {
	go room.listenActions()
}

func (room *ChatRoom) listenActions() {
	for {
		select {
		case info := <-room.additionQueue:
			peer := NewPeer(info.conn, info.addr, info.name)
			room.registerPeer(peer)
		case info := <-room.removalQueue:
			room.unregisterPeer(info.addr, info.name)
		}
	}
}

func (room *ChatRoom) registerPeer(peer *Peer) {
	room.peers = append(room.peers, peer)
	go func() {
		for msg := range peer.Incoming {
			room.Incoming <- msg
		}
	}()
}

func (room *ChatRoom) unregisterPeer(addr string, name string) {
	iToRemove := -1
	for i, v := range room.peers {
		if v.Addr == addr && v.Name == name {
			iToRemove = i
		}
	}
	if i := iToRemove; i != -1 {
		room.peers = append(room.peers[:i], room.peers[i+1:]...)
	}
}

type peerRemoval struct {
	addr string
	name string
}

type peerAddition struct {
	conn net.Conn
	addr string
	name string
}
