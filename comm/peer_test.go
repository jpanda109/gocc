package comm

import (
	"bytes"
	"testing"
)

func TestPeer(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewPeer(buf, buf, "localhost:8080", "anon")
	var err error
	var msg *MsgGob
	err = p.Send(&MsgGob{Public, "asdf"})
	if err != nil {
		t.Error("peer.Send(*MsgGob) wrongly gave error")
	}
	msg, err = p.Receive()
	if err != nil {
		t.Error("peer.Receive() wrongly received error")
	}
	if *msg != *(&MsgGob{Public, "asdf"}) {
		t.Error("peer.Receive() received corrupt data")
	}
}
