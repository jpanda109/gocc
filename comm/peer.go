package comm

import (
	"encoding/gob"
	"errors"
	"net"
)

var curID int

// idIncrementer uses closure to return a unique ID every time
func idIncrementer() int {
	curID++
	return curID
}

// NewPeer returns a new peer
func NewPeer(conn net.Conn, addr string, name string) *Peer {
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)
	newPeer := &Peer{
		idIncrementer(),
		addr,
		name,
		make(chan *MsgGob),
		make(chan *MsgGob),
		decoder,
		encoder,
	}
	return newPeer
}

// Peer represents a peer server
// ID is a unique id for each peer
// Addr is the remote address of the peer
// outgoing is a channel of gobs used to send data to the peer
// incoming is a channel of gobs user to receive data from the peer
// decoder sends a gob over TCP to the peer
// encoder receives a gob over TCP from the peer
type Peer struct {
	ID       int
	Addr     string
	Name     string
	outgoing chan *MsgGob
	incoming chan *MsgGob
	decoder  *gob.Decoder
	encoder  *gob.Encoder
}

// Send sends message to peer, returns error if closed
func (p *Peer) Send(msg *MsgGob) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				err = x
			}
		}
	}()
	p.outgoing <- msg
	return nil
}

// Receive returns message from peer, blocks until available, err if closed
func (p *Peer) Receive() (*MsgGob, error) {
	msg, ok := <-p.incoming
	if !ok {
		return nil, errors.New("")
	}
	return msg, nil
}

func (p *Peer) String() string {
	return p.Addr + "," + p.Name
}

// Start begins listening functions
func (p *Peer) Start() {
	go p.beginRead()
	go p.beginWrite()
}

func (p *Peer) quit() {
	close(p.outgoing)
	close(p.incoming)
}

func (p *Peer) beginRead() {
	for {
		msg := &MsgGob{}
		err := p.decoder.Decode(msg)
		if err != nil {
			p.quit()
			return
		}
		p.incoming <- msg
	}
}

func (p *Peer) beginWrite() {
	for msg := range p.outgoing {
		p.encoder.Encode(msg)
	}
}
