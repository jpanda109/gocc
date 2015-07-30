package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jpanda109/gocc/input"
	"github.com/nsf/termbox-go"
)

const msglen = 256

func setLogger() {
	fo, _ := os.Create("log.txt")
	log.SetOutput(fo)
}

func startApp(port string, debug bool, connect string, name string) {
	setLogger()
	listenerAddr := ":" + port
	if debug {
		listenerAddr = "localhost" + listenerAddr
	}

	controller := input.NewController(listenerAddr, name)
	if connect != "" {
		err := controller.Connect(connect)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	termbox.Init()
	defer termbox.Close()
	wg := controller.Start()
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "8080",
			Usage: "Port listen to",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Sets server to localhost only",
		},
		cli.StringFlag{
			Name:  "connect, c",
			Value: "",
			Usage: "Address to connect to",
		},
		cli.StringFlag{
			Name:  "name, n",
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
