package comm

import "sync"

// NewChatRoom creates and returns a pointer to chat room
func NewChatRoom() *ChatRoom {
	room := &ChatRoom{
		make(chan *Message),
		make(chan string),
		&sync.Mutex{},
		[]*Peer{},
	}
	return room
}

// ChatRoom sends and receives messages to peers
type ChatRoom struct {
	incoming chan *Message
	outgoing chan string
	peerLock *sync.Mutex
	peers    []*Peer
}

// Broadcast sends message to all peers
func (room *ChatRoom) Broadcast(msg string) {
	for _, peer := range room.peers {
		peer.Send(msg)
	}
}

// Receive returns the next message from any peer
func (room *ChatRoom) Receive() *Message {
	return <-room.incoming
}

// AddPeer adds peer to chat room with given info
func (room *ChatRoom) AddPeer(peer *Peer) {
	room.peerLock.Lock()
	defer room.peerLock.Unlock()
	room.peers = append(room.peers, peer)
	go func() {
		for {
			msg, err := peer.Receive()
			if err != nil {
				room.RemovePeer(peer.Addr, peer.Name)
				break
			}
			room.incoming <- &Message{peer, msg}
		}
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
