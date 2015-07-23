package comm

import (
	"net"
	"sync"
)

// NewChatRoom creates and returns a pointer to chat room
func NewChatRoom() *ChatRoom {
	room := &ChatRoom{
		make(chan string),
		make(chan string),
		&sync.Mutex{},
		[]*Peer{},
	}
	return room
}

// ChatRoom sends and receives messages to peers
type ChatRoom struct {
	Incoming chan string
	Outgoing chan string
	peerLock *sync.Mutex
	peers    []*Peer
}

// AddPeer adds peer to chat room with given info
func (room *ChatRoom) AddPeer(conn net.Conn, addr string, name string) {
	room.peerLock.Lock()
	defer room.peerLock.Unlock()
	peer := NewPeer(conn, addr, name)
	room.peers = append(room.peers, peer)
	go func() {
		for msg := range peer.Incoming {
			room.Incoming <- msg
		}
		// peer has disconnected at this point
		room.RemovePeer(peer.Addr, peer.Name)
	}()
}

// RemovePeer removes peer from chat room with given info
func (room *ChatRoom) RemovePeer(addr string, name string) {
	room.peerLock.Lock()
	defer room.peerLock.Unlock()
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
