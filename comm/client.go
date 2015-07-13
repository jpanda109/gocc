package comm

import (
	"bufio"
	"net"
	"strings"
)

// NewClient initializes client to 0 state
func NewClient(conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	newClient := &Client{
		"",
		make(chan string),
		make(chan string),
		make(chan bool),
		reader,
		writer,
	}
	newClient.Start()
	return newClient
}

// Client allows for bidirectional communication
type Client struct {
	Name     string
	Outgoing chan string
	Incoming chan string
	Quit     chan bool
	reader   *bufio.Reader
	writer   *bufio.Writer
}

// Start read and write functions
func (c *Client) Start() {
	go c.beginRead()
	go c.beginWrite()
	<-c.Quit
}

// BeginRead reads input from client connection and streams to client's
// read channel
func (c *Client) beginRead() {
	for {
		msg, _ := c.reader.ReadString('\n')
		msg = strings.Trim(msg, "\n")
		c.Incoming <- msg
	}
}

// BeginWrite takes input from client's writing channel and write to
// client's connection
func (c *Client) beginWrite() {
	for msg := range c.Outgoing {
		c.writer.WriteString(msg)
		c.writer.Flush()
	}
}
