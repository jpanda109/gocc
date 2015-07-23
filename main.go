package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jpanda109/gocc/comm"
)

const msglen = 256

func startApp(port string, debug bool, connect string, name string) {
	listenerAddr := ":" + port
	if debug {
		listenerAddr = "localhost" + listenerAddr
	}
	handler := comm.NewConnHandler(listenerAddr, name)
	if connect != "" {
		handler.Dial(connect)
	}
	handler.Listen()
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "Port listen to",
		},
		cli.BoolTFlag{
			Name:  "debug",
			Usage: "Sets server to localhost only",
		},
		cli.StringFlag{
			Name:  "connect",
			Value: "",
			Usage: "Address to connect to",
		},
		cli.StringFlag{
			Name:  "name",
			Value: "anon",
			Usage: "Name with which to identify",
		},
	}
	app.Action = func(c *cli.Context) {
		startApp(
			c.String("port"),
			c.BoolT("debug"),
			c.String("connect"),
			c.String("name"),
		)
	}
	app.Run(os.Args)
}
