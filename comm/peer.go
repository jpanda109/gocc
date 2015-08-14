package comm

import (
	"encoding/gob"
	"errors"
	"io"
)

var curID int

// idIncrementer uses closure to return a unique ID every time
func idIncrementer() int {
	curID++
	return curID
}

// NewPeer returns a new peer
func NewPeer(r io.Reader, w io.Writer, addr string, name string) *Peer {
	decoder := gob.NewDecoder(r)
	encoder := gob.NewEncoder(w)
	newPeer := &Peer{
		idIncrementer(),
		addr,
		name,
		decoder,
		encoder,
	}
	return newPeer
}

// Peer represents a peer server
// ID is a unique id for each peer
// Addr is the remote address of the peer
// decoder sends a gob over TCP to the peer
// encoder receives a gob over TCP from the peer
type Peer struct {
	ID      int
	Addr    string
	Name    string
	decoder *gob.Decoder
	encoder *gob.Encoder
}

// Send sends message to peer, returns error if closed
func (p *Peer) Send(msg *MsgGob) error {
	err := p.encoder.Encode(msg)
	return err
}

// Receive returns message from peer, blocks until available, err if closed
func (p *Peer) Receive() (*MsgGob, error) {
	msg := &MsgGob{}
	err := p.decoder.Decode(msg)
	if err != nil {
		return nil, errors.New("")
	}
	return msg, nil
}

func (p *Peer) String() string {
	return p.Addr + "," + p.Name
}
