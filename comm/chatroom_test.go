package comm

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	ln, _ := net.Listen("tcp", "localhost:8080")
	chatroom := NewChatRoom()
	connCh := make(chan net.Conn)
	go listen2(ln, connCh)
	go func() {
		for c := range connCh {
			chatroom.NewConnections <- c
		}
	}()
	var readers []*bufio.Reader
	var cConns []net.Conn
	for i := 0; i < 2; i++ {
		cConn, _ := net.Dial("tcp", "localhost:8080")
		cConns = append(cConns, cConn)
		cWriter := bufio.NewWriter(cConn)
		cWriter.WriteString("NAME" + strconv.Itoa(i) + "\n")
		cWriter.Flush()
		readers = append(readers, bufio.NewReader(cConn))
	}
	time.Sleep(500 * time.Millisecond) // hack please change
	chatroom.Outgoing <- "Hello World\n"
	for _, r := range readers {
		if s, _ := r.ReadString('\n'); strings.Trim(s, "\n") != "Hello World" {
			t.Error("error")
		}
	}
	ln.Close()
}

func listen2(ln net.Listener, connCh chan net.Conn) {
	for {
		conn, _ := ln.Accept()
		go handleConn(conn, connCh)
	}
}

func handleConn(conn net.Conn, connCh chan net.Conn) {
	connCh <- conn
}
