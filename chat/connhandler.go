package chat

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// NewConnHandler creates a new conn handler
func NewConnHandler(addr string, name string, chatroom Room) *ConnHandler {
	handler := &ConnHandler{
		addr,
		name,
		make(chan Peer),
		chatroom,
		make(chan string),
		[]string{},
	}
	return handler
}

// ConnHandler listens for new connections and connets to new ones
type ConnHandler struct {
	addr       string
	name       string
	newPeers   chan Peer
	chatroom   Room
	whitelists chan string
	whitelist  []string
}

// String returns string repr of conn handler
func (handler *ConnHandler) String() string {
	return handler.addr + "," + handler.name
}

// Peer returns the next peer that connects
func (handler *ConnHandler) Peer() Peer {
	return <-handler.newPeers
}

// Whitelist adds an ip address to the whitelist
func (handler *ConnHandler) Whitelist(addr string) {
	handler.whitelists <- addr
}

// Dial connects to server at address
func (handler *ConnHandler) Dial(addr string) ([]Peer, error) {
	var peers []Peer
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, "\n")
	addrs := strings.Split(line, ";")
	info := strings.Split(addrs[0], ",")
	writer.WriteString(handler.String() + "\n")
	writer.Flush()
	peers = append(peers, NewPeer(conn, conn, info[0], info[1]))
	for _, a := range addrs[1:] {
		info := strings.Split(a, ",")
		c, _ := net.Dial("tcp", info[0])
		writer = bufio.NewWriter(c)
		writer.WriteString(handler.String() + "\n")
		writer.Flush()
		reader = bufio.NewReader(c)
		line, _ = reader.ReadString('\n')
		peers = append(peers, NewPeer(c, c, info[0], info[1]))
	}
	return peers, nil
}

// Listen begins listener
func (handler *ConnHandler) Listen() {
	go handler.listenConns()
}

func (handler *ConnHandler) handleWhitelists() {
	for addr := range handler.whitelists {
		handler.whitelist = append(handler.whitelist, addr)
	}
}

func (handler *ConnHandler) listenConns() {
	ln, err := net.Listen("tcp", handler.addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err == nil {
			writer := bufio.NewWriter(conn)
			reader := bufio.NewReader(conn)
			writer.WriteString(handler.String())
			for _, p := range handler.chatroom.Peers() {
				writer.WriteString(";")
				writer.WriteString(p.String())
			}
			writer.WriteString("\n")
			writer.Flush()
			line, _ := reader.ReadString('\n')
			line = strings.Trim(line, "\n")
			info := strings.Split(line, ",")
			handler.newPeers <- NewPeer(conn, conn, info[0], info[1])
		}
	}
}
