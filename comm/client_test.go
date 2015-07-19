package comm

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"
)

var (
	ln net.Listener
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestClientRead(t *testing.T) {
	ln, _ := net.Listen("tcp", "localhost:8080")
	var sConn net.Conn
	go func() {
		sConn, _ = ln.Accept()
	}()
	cConn, _ := net.Dial("tcp", "localhost:8080")
	client := NewClient(sConn)
	cWriter := bufio.NewWriter(cConn)
	// cReader := bufio.NewReader(cConn)
	go func() {
		cWriter.WriteString("HI\n")
		cWriter.WriteString("BYE\n")
		cWriter.Flush()
	}()
	if hi := <-client.Incoming; hi != "HI" {
		t.Error("error")
	}
	if bye := <-client.Incoming; bye != "BYE" {
		t.Error("error")
	}
	cConn.Close()
	sConn.Close()
	ln.Close()
}

func TestClientWrite(t *testing.T) {
	fmt.Println("testclientwrite")
	ln, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		fmt.Println(err)
	}
	connCh := make(chan net.Conn)
	go listen(ln, connCh)
	cConn, _ := net.Dial("tcp", "localhost:8090")
	if err != nil {
		fmt.Println(err)
	}
	sConn := <-connCh
	if sConn == nil {
		fmt.Println("sConn is nil")
	}
	fmt.Println(sConn.LocalAddr())
	client := NewClient(sConn)
	cWriter := bufio.NewWriter(cConn)
	// cReader := bufio.NewReader(cConn)
	go func() {
		cWriter.WriteString("HI\n")
		cWriter.WriteString("BYE\n")
		cWriter.Flush()
	}()
	if hi := <-client.Incoming; hi != "HI" {
		t.Error("error")
	}
	if bye := <-client.Incoming; bye != "BYE" {
		t.Error("error")
	}
	cConn.Close()
	sConn.Close()
	ln.Close()
	// fmt.Println("test write")
	// ln, err := net.Listen("tcp", "localhost:8090")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var sConn net.Conn
	// go func() {
	// 	sConn, err = ln.Accept()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	// cConn, _ := net.Dial("tcp", "localhost:8090")
	// client := NewClient(sConn)
	// cReader := bufio.NewReader(cConn)
	// go func() {
	// 	client.Outgoing <- "HI\n"
	// 	client.Outgoing <- "BYE\n"
	// }()
	// if hi, _ := cReader.ReadString('\n'); hi != "HI\n" {
	// 	t.Error("error")
	// }
	// if bye, _ := cReader.ReadString('\n'); bye != "BYE\n" {
	// 	t.Error("error")
	// }
	// cConn.Close()
	// sConn.Close()
	// ln.Close()
}

func listen(ln net.Listener, connCh chan net.Conn) {
	for {
		fmt.Println("accept")
		conn, _ := ln.Accept()
		connCh <- conn
	}
}
