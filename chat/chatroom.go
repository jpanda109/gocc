package chat

import "sync"

// Room defines an interface which can receive and broadcast messages to
// a number of peers
type Room interface {
	Broadcast(msg string)
	Receive() *Message
	AddPeer(peer Peer)
	RemovePeer(peer Peer)
	Peers() []Peer
}

// NewRoom creates and returns a pointer to chat room
func NewRoom() Room {
	room := &chatRoom{
		make(chan *Message),
		make(chan string),
		&sync.Mutex{},
		[]Peer{},
	}
	return room
}

// Room sends and receives messages to peers
// incoming is a channel of messages from each peer
// peerLock handles atomicity of adding and removing peers
// peers is a list of peers in the chat room
type chatRoom struct {
	incoming   chan *Message
	broadcasts chan string
	peerLock   *sync.Mutex
	peers      []Peer
}

// Broadcast sends message to all peers
func (room *chatRoom) Broadcast(msg string) {
	room.broadcasts <- msg
}

// Receive returns the next message from any peer
func (room *chatRoom) Receive() *Message {
	return <-room.incoming
}

// AddPeer adds peer to chat room with given info
func (room *chatRoom) AddPeer(peer Peer) {
	room.peerLock.Lock()
	defer room.peerLock.Unlock()
	room.peers = append(room.peers, peer)
	go func() {
		for {
			msg, err := peer.Receive()
			if err != nil {
				room.RemovePeer(peer)
				break
			}
			room.incoming <- &Message{peer.ID(), peer.Name(), msg}
		}
	}()
	go func() {
		for {
			msg := <-room.broadcasts
			for _, peer := range room.peers {
				peer.Send(&MsgGob{Broadcast, msg})
			}
		}
	}()
}

// RemovePeer removes peer from chat room with given info
func (room *chatRoom) RemovePeer(peer Peer) {
	room.peerLock.Lock()
	defer room.peerLock.Unlock()
	iToRemove := -1
	for i, v := range room.peers {
		if v == peer {
			iToRemove = i
		}
	}
	if i := iToRemove; i != -1 {
		room.peers = append(room.peers[:i], room.peers[i+1:]...)
	}
}

// Peers returns a list of peers
func (room *chatRoom) Peers() []Peer {
	return room.peers
}
