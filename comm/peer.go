package comm

import (
	"bufio"
	"net"
)

// NewPeer returns a new peer
func NewPeer(conn net.Conn, addr string, name string) *Peer {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	newPeer := &Peer{
		addr,
		name,
		make(chan string),
		make(chan string),
		make(chan bool),
		reader,
		writer,
	}
	return newPeer
}

// Peer represents a peer server
type Peer struct {
	Addr       string
	Name       string
	Outgoing   chan string
	Incoming   chan string
	Disconnect chan bool
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func (p *Peer) String() string {
	return p.Addr + "," + p.Name
}
