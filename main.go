package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jpanda109/gocc/app"
	"github.com/jpanda109/gocc/chat"
	"github.com/jpanda109/gocc/config"
	"github.com/jpanda109/gocc/friends"
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
	manager := &app.Manager{}
	chatApp := &chat.App{
		Manager: manager,
		Addr:    listenerAddr,
		Name:    name,
	}
	termbox.Init()
	defer termbox.Close()
	wg, _ := manager.Start(chatApp)
	wg.Wait()
	// listenerAddr := ":" + port
	// if debug {
	// 	listenerAddr = "localhost" + listenerAddr
	// }
	// setLogger(debug)
	// controller := chat.NewController(listenerAddr, name)
	// if connect != "" {
	// 	err := controller.Connect(connect)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }
	// termbox.Init()
	// defer termbox.Close()
	// wg := controller.Start()
	// wg.Wait()
}

func start(screen app.Screen) {
	termbox.Init()
	defer termbox.Close()
	manager := &app.Manager{}
	screen.SetManager(manager)
	wg, _ := manager.Start(screen)
	wg.Wait()
}

// main simply handles command line flag processing and then calls
// the startApp function with the correct parameters
// It will pass in the appropriate flags depending on both the config file
// and command line flags, with command line flags having higher precedence
func main() {
	config.Init()
	app := cli.NewApp()
	config.SetDebug(true)
	setLogger(true)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: config.Port(),
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
			Value: config.Name(),
			Usage: "Name with which to identify",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "friends",
			Aliases: []string{"f"},
			Usage:   "display friends",
			Action: func(c *cli.Context) {
				start(&friends.FriendApp{})
			},
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
