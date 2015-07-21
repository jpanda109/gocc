package comm

import (
	"bufio"
	"net"
	"strings"
	"sync"
)

// NewClient initializes client to 0 state, blocks until connection sends name
func NewClient(conn net.Conn, name string) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	newClient := &Client{
		name,
		make(chan string),
		make(chan string),
		make(chan bool),
		&sync.WaitGroup{},
		reader,
		writer,
	}
	newClient.start()
	return newClient
}

// Client allows for bidirectional communication
type Client struct {
	Name     string
	Outgoing chan string
	Incoming chan string
	Quit     chan bool
	wg       *sync.WaitGroup
	reader   *bufio.Reader
	writer   *bufio.Writer
}

// Start read and write functions
func (c *Client) start() {
	c.wg.Add(2)
	go c.beginRead()
	go c.beginWrite()
}

// Stop stuff
func (c *Client) Stop() {
	close(c.Quit)
	c.wg.Wait()
}

// BeginRead reads input from client connection and streams to client's
// read channel
func (c *Client) beginRead() {
	defer c.wg.Done()
	for {
		msg, err := c.reader.ReadString('\n')
		if err != nil {
			c.Stop()
			return
		}
		msg = strings.Trim(msg, "\n")
		c.Incoming <- msg
	}
}

// BeginWrite takes input from client's writing channel and write to
// client's connection
func (c *Client) beginWrite() {
	defer c.wg.Done()
	for {
		select {
		case <-c.Quit:
			return
		case msg := <-c.Outgoing:
			if !strings.HasSuffix(msg, "\n") {
				msg = msg + "\n"
			}
			c.writer.WriteString(msg)
			c.writer.Flush()
		}
	}
}
