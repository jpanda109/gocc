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
		[]string{addr + "," + name},
		make(chan *Client),
	}
	handler.Listen()
	return handler
}

// ConnHandler listens for new connections and connets to new ones
type ConnHandler struct {
	addr       string
	addrs      []string
	NewClients chan *Client
}

// Listen begins listener
func (handler *ConnHandler) Listen() {
	go handler.listenConns()
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
	handler.addrs = append(handler.addrs, addrs...)
	writer.WriteString(handler.addrs[0] + "\n")
	writer.Flush()
	handler.NewClients <- NewClient(conn, strings.Split(addrs[0], ",")[1])
	for _, a := range addrs[1:] {
		info := strings.Split(a, ",")
		c, _ := net.Dial("tcp", info[0])
		writer = bufio.NewWriter(c)
		writer.WriteString(handler.addrs[0] + "\n")
		writer.Flush()
		reader = bufio.NewReader(c)
		line, _ = reader.ReadString('\n')
		handler.NewClients <- NewClient(c, info[1])
	}
	fmt.Println(handler.addrs)
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
			for i, a := range handler.addrs {
				writer.WriteString(a)
				if i < len(handler.addrs)-1 {
					writer.WriteString(";")
				}
			}
			writer.WriteString("\n")
			writer.Flush()
			line, _ := reader.ReadString('\n')
			line = strings.Trim(line, "\n")
			info := strings.Split(line, ",")
			name := info[1]
			handler.NewClients <- NewClient(conn, name)
			handler.addrs = append(handler.addrs, line)
			fmt.Println(handler.addrs)
		}
	}
}
