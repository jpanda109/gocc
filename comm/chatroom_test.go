package comm

import "testing"

type mockPeer struct {
	id   int
	name string
}

func (p *mockPeer) ID() int {
	return p.id
}

func (p *mockPeer) Name() string {
	return p.name
}

func (p *mockPeer) Send(msg *MsgGob) error {
	return nil
}

func (p *mockPeer) Receive() (*MsgGob, error) {
	return nil, nil
}

func (p *mockPeer) String() string {
	return p.name
}

var peers = []Peer{
	&mockPeer{0, "n0"},
	&mockPeer{1, "n1"},
	&mockPeer{2, "n2"},
	&mockPeer{3, "n3"},
}

func seededChatRoom() ChatRoom {
	room := NewChatRoom()
	for _, p := range peers {
		room.AddPeer(p)
	}
	return room
}

func TestChatroomAddPeer(t *testing.T) {
	room := seededChatRoom()
	if len(room.Peers()) != len(peers) {
		t.Error("Chatroom didn't add proper number of peers")
	}
	for i, p := range room.Peers() {
		if p != peers[i] {
			t.Error("add multiple peers failed")
		}
	}
}

func TestChatroomRemovePeer(t *testing.T) {
	room := seededChatRoom()
	room.RemovePeer(peers[0])
	if len(room.Peers()) != len(peers)-1 {
		t.Error("Chatroom didn't remove any peers (first element)")
	}
	for i, p := range room.Peers() {
		if p != peers[i+1] {
			t.Error("Chatroom didn't remove correct peer (first element)")
		}
	}
	room = seededChatRoom()
	room.RemovePeer(peers[len(peers)-1])
	if len(room.Peers()) != len(peers)-1 {
		t.Error("Chatroom didn't remove any peers (last element)")
	}
	for i, p := range room.Peers() {
		if p != peers[i] {
			t.Error("Chatroom didn't remove correct peer (last element)")
		}
	}
	room = seededChatRoom()
	room.RemovePeer(peers[1])
	if len(room.Peers()) != len(peers)-1 {
		t.Error("Chatroom didn't remove any peers (non-edge element)")
	}
	for i, p := range append(peers[:1], peers[2:]...) {
		if p != room.Peers()[i] {
			t.Error("Chatroom didnt remove correct peer (non-edge element)")
		}
	}
}
