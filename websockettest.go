package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

type T struct {
	Msg   string
	Count int
}

func main(){
	conn, err := websocket.Dial("wss://echo.websocket.org", "", "http://localhost")
	if err != nil {
		log.Fatalf("error dialing websocket; %v", err)
	}
	defer conn.Close()
	if err = websocket.JSON.Send(conn, "test123"); err != nil {
		log.Fatalf("error sending request; %v", err)
	}
	var message string
	if err = websocket.JSON.Receive(conn, &message); err != nil {
		log.Fatalf("error receiving response; %v", err)
	}
	fmt.Printf("%+v\n", message)
}
