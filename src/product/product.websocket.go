package product

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {
	fmt.Println("new websocket connection established")
	// Keep looping until a EOF message is received from client
	for {
		var msg message
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Println(err)
			break
		}
		fmt.Printf("received message %s\n", msg.Data)
	}
	fmt.Println("closing the websocket")
	defer ws.Close()

	// test it with the following 2 lines of javascript
	/*
		let ws = new WebSocket("ws://localhost:5000/websocket")
		ws.send(JSON.stringify({data:"test message from browser", type:"test"}))
	*/
}
