package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jpanda109/gocc/input"
	"github.com/nsf/termbox-go"
)

const msglen = 256

func startApp(port string, debug bool, connect string, name string) {
	listenerAddr := ":" + port
	if debug {
		listenerAddr = "localhost" + listenerAddr
	}
	// chatroom := comm.NewChatRoom()
	// handler := comm.NewConnHandler(listenerAddr, name, chatroom)
	// if connect != "" {
	// 	peers := handler.Dial(connect)
	// 	for _, peer := range peers {
	// 		chatroom.AddPeer(peer)
	// 	}
	// }
	// handler.Listen()
	// go func() {
	// 	for {
	// 		peer := handler.GetPeer()
	// 		chatroom.AddPeer(peer)
	// 	}
	// }()
	// controller := input.Handler{
	// 	handler,
	// 	chatroom,
	// 	view.NewChatWindow(),
	// 	make([]rune, 0),
	// }
	termbox.Init()
	defer termbox.Close()
	controller := input.NewHandler(listenerAddr, name)
	if connect != "" {
		controller.Connect(connect)
	}
	wg := controller.Start()
	wg.Wait()

	// reader := bufio.NewReader(os.Stdin)
	// go func() {
	// 	for {
	// 		msg := chatroom.Receive()
	// 		fmt.Println(msg)
	// 	}
	// }()
	// for {
	// 	msg, _ := reader.ReadString('\n')
	// 	chatroom.Broadcast(msg)
	// }
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "8080",
			Usage: "Port listen to",
		},
		cli.BoolTFlag{
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
