package comm

import (
	"bufio"
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
	Outgoing chan string
	Incoming chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (p *Peer) String() string {
	return p.Addr + "," + p.Name
}

func (p *Peer) start() {
	go p.beginRead()
	go p.beginWrite()
}

func (p *Peer) quit() {
	close(p.Outgoing)
	close(p.Incoming)
}

func (p *Peer) beginRead() {
	for {
		msg, err := p.reader.ReadString('\n')
		if err != nil {
			p.quit()
			return
		}
		msg = strings.Trim(msg, "\n")
		p.Incoming <- msg
	}
}

func (p *Peer) beginWrite() {
	for msg := range p.Outgoing {
		if !strings.HasSuffix(msg, "\n") {
			msg = msg + "\n"
		}
		p.writer.WriteString(msg)
		p.writer.Flush()
	}
}
