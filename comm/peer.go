package comm

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
)

// MsgGob is a gob which used to send necessary data through tcp sockets
// Trusted information is not contained due to possibility of unreliable clients
// action is the type of message being sent
// body is the body of the message being sent
type MsgGob struct {
	Action Action
	Body   string
}

// Action represents the type of action that the message represents
type Action int

const (
	// Public indicates the message is broadcasted to all
	Public Action = iota
	// Private indicates the message was sent to one peer
	Private
	// Whitelist indicates a modification to the chatroom whitelist
	Whitelist
)

// Message contains information about a message received from a peer
// This struct, unlike MsgGob, is meant to be passed internally within the
// application with trusted information
// Sender is a pointer to the local peer instance which sent it
// info contains the message infomation such as the action type and body
type Message struct {
	SenderID   int
	SenderName string
	Info       *MsgGob
}

// String returns a human readable form of a Message
func (msg *Message) String() string {
	return fmt.Sprintf("%s (%v) > %s", msg.SenderName, msg.Info.Action,
		msg.Info.Body)
}

var curID int

// idIncrementer uses closure to return a unique ID every time
func idIncrementer() int {
	curID++
	return curID
}

// NewPeer returns a new peer
func NewPeer(r io.Reader, w io.Writer, addr string, name string) Peer {
	decoder := gob.NewDecoder(r)
	encoder := gob.NewEncoder(w)
	newPeer := &peer{
		idIncrementer(),
		addr,
		name,
		decoder,
		encoder,
	}
	return newPeer
}

// Peer defines an interface which can send and receive messages from a resource
type Peer interface {
	ID() int
	Name() string
	Send(msg *MsgGob) error
	Receive() (*MsgGob, error)
	String() string
}

// peer represents a peer server
// ID is a unique id for each peer
// Addr is the remote address of the peer
// decoder sends a gob over TCP to the peer
// encoder receives a gob over TCP from the peer
type peer struct {
	id      int
	addr    string
	name    string
	decoder *gob.Decoder
	encoder *gob.Encoder
}

// ID is an accessor method for id
func (p *peer) ID() int {
	return p.id
}

// Name is an accessor method for name
func (p *peer) Name() string {
	return p.name
}

// Send sends message to peer, returns error if closed
func (p *peer) Send(msg *MsgGob) error {
	err := p.encoder.Encode(msg)
	return err
}

// Receive returns message from peer, blocks until available, err if closed
func (p *peer) Receive() (*MsgGob, error) {
	msg := &MsgGob{}
	err := p.decoder.Decode(msg)
	if err != nil {
		return nil, errors.New("")
	}
	return msg, nil
}

func (p *peer) String() string {
	return p.addr + "," + p.name
}
