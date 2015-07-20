package notmain

import (
	"fmt"
	"net"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jpanda109/gocc/comm"
)

const msglen = 256

func startApp(port string, debug bool) {
	listenerAddr := ":" + port
	if debug {
		listenerAddr = "localhost" + listenerAddr
	}
	ln, err := net.Listen("tcp", listenerAddr)
	if err != nil {
		fmt.Println("Error in setting up tcp listener")
		os.Exit(-1)
	}
	fmt.Printf("Listening at %s\n", ln.Addr())
	chatroom := comm.NewChatRoom()
	for {
		conn, err := ln.Accept()
		if err != nil {
		}
		chatroom.NewConnections <- conn
	}
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "Port to connect or listen to",
		},
		cli.BoolTFlag{
			Name:  "debug",
			Usage: "Sets server to localhost only",
		},
	}
	app.Action = func(c *cli.Context) {
		startApp(
			c.String("port"),
			c.BoolT("debug"),
		)
	}
	app.Run(os.Args)
}
