package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jpanda109/gocc/config"
	"github.com/jpanda109/gocc/input"
	"github.com/nsf/termbox-go"
)

// setLogger either logs to a log file or discards log statements depending
// on the passed in boolean
func setLogger(debug bool) {
	if debug {
		fo, _ := os.Create("log.txt")
		log.SetOutput(fo)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

// startApp creates the controller and runs the application
func startApp(port string, debug bool, connect string, name string) {
	listenerAddr := ":" + port
	if debug {
		listenerAddr = "localhost" + listenerAddr
	}
	setLogger(debug)
	controller := input.NewController(listenerAddr, name)
	if connect != "" {
		err := controller.Connect(connect)
		if err != nil {
			log.Println(err)
			return
		}
	}
	termbox.Init()
	defer termbox.Close()
	wg := controller.Start()
	wg.Wait()
}

// main simply handles command line flag processing and then calls
// the startApp function with the correct parameters
// It will pass in the appropriate flags depending on both the config file
// and command line flags, with command line flags having higher precedence
func main() {
	config.Init()
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
