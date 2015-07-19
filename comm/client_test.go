package comm

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

//
// func TestClientRead(t *testing.T) {
// 	ln, _ := net.Listen("tcp", "localhost:8010")
// 	var sConn net.Conn
// 	go func() {
// 		sConn, _ = ln.Accept()
// 	}()
// 	cConn, _ := net.Dial("tcp", "localhost:8010")
// 	cWriter := bufio.NewWriter(cConn)
// 	// cReader := bufio.NewReader(cConn)
// 	go func() {
// 		cWriter.WriteString("NAME\n")
// 		cWriter.WriteString("HI\n")
// 		cWriter.WriteString("BYE\n")
// 		cWriter.Flush()
// 	}()
// 	client := NewClient(sConn)
// 	if hi := <-client.Incoming; hi != "HI" {
// 		t.Error("error")
// 	}
// 	if bye := <-client.Incoming; bye != "BYE" {
// 		t.Error("error")
// 	}
// 	cConn.Close()
// 	sConn.Close()
// 	ln.Close()
// }
//
// func TestClientWrite(t *testing.T) {
// 	ln, _ := net.Listen("tcp", "localhost:8090")
// 	connCh := make(chan net.Conn)
// 	go listen(ln, connCh)
// 	cConn, _ := net.Dial("tcp", "localhost:8090")
// 	sConn := <-connCh
// 	cWriter := bufio.NewWriter(cConn)
// 	cReader := bufio.NewReader(cConn)
// 	go func() {
// 		cWriter.WriteString("NAME\n")
// 		cWriter.Flush()
// 	}()
// 	client := NewClient(sConn)
// 	go func() {
// 		client.Outgoing <- "HI\n"
// 		client.Outgoing <- "BYE\n"
// 	}()
// 	if hi, _ := cReader.ReadString('\n'); hi != "HI\n" {
// 		t.Error("error")
// 	}
// 	if bye, _ := cReader.ReadString('\n'); bye != "BYE\n" {
// 		t.Error("error")
// 	}
// 	cConn.Close()
// 	sConn.Close()
// 	ln.Close()
// }
//
// func listen(ln net.Listener, connCh chan net.Conn) {
// 	for {
// 		conn, _ := ln.Accept()
// 		connCh <- conn
// 	}
// }
