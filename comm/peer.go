package comm

import (
	"bufio"
	"errors"
	"net"
	"strings"
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
		reader,
		writer,
	}
	newPeer.start()
	return newPeer
}

// Peer represents a peer server
type Peer struct {
	Addr     string
	Name     string
	outgoing chan string
	incoming chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

// Send sends message to peer, returns error if closed
func (p *Peer) Send(msg string) (err error) {
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
func (p *Peer) Receive() (string, error) {
	msg, ok := <-p.incoming
	if !ok {
		return "", errors.New("")
	}
	return msg, nil
}

func (p *Peer) String() string {
	return p.Addr + "," + p.Name
}

func (p *Peer) start() {
	go p.beginRead()
	go p.beginWrite()
}

func (p *Peer) quit() {
	close(p.outgoing)
	close(p.incoming)
}

func (p *Peer) beginRead() {
	for {
		msg, err := p.reader.ReadString('\n')
		if err != nil {
			p.quit()
			return
		}
		msg = strings.Trim(msg, "\n")
		p.incoming <- msg
	}
}

func (p *Peer) beginWrite() {
	for msg := range p.outgoing {
		if !strings.HasSuffix(msg, "\n") {
			msg = msg + "\n"
		}
		p.writer.WriteString(msg)
		p.writer.Flush()
	}
}
