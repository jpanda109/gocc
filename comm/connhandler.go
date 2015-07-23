package comm

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// NewConnHandler creates a new conn handler
func NewConnHandler(addr string, name string) *ConnHandler {
	handler := &ConnHandler{
		addr,
		name,
		NewChatRoom(),
	}
	return handler
}

// ConnHandler listens for new connections and connets to new ones
type ConnHandler struct {
	addr     string
	name     string
	chatroom *ChatRoom
}

// String returns string repr of conn handler
func (handler *ConnHandler) String() string {
	return handler.addr + "," + handler.name
}

// Listen begins listener
func (handler *ConnHandler) Listen() {
	handler.listenConns()
}

// Dial connects to server at address
func (handler *ConnHandler) Dial(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, "\n")
	addrs := strings.Split(line, ";")
	for _, a := range addrs {
		info := strings.Split(a, ",")
		handler.chatroom.AddPeer(conn, info[0], info[1])
	}
	writer.WriteString(handler.String() + "\n")
	writer.Flush()
	for _, a := range addrs[1:] {
		info := strings.Split(a, ",")
		c, _ := net.Dial("tcp", info[0])
		writer = bufio.NewWriter(c)
		writer.WriteString(handler.String() + "\n")
		writer.Flush()
		reader = bufio.NewReader(c)
		line, _ = reader.ReadString('\n')
	}
	fmt.Println(handler.chatroom.peers)
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
			for _, p := range handler.chatroom.peers {
				writer.WriteString(";")
				writer.WriteString(p.String())
			}
			writer.WriteString("\n")
			writer.Flush()
			line, _ := reader.ReadString('\n')
			line = strings.Trim(line, "\n")
			info := strings.Split(line, ",")
			handler.chatroom.AddPeer(conn, info[0], info[1])
		}
		fmt.Println(handler.chatroom.peers)
	}
}
