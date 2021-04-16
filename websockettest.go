package main

import (
	"fmt"
	//"golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type URL struct {
	Schema string
	Host string
	Path string
}

func main(){
	/*conn, err := websocket.Dial("wss://echo.websocket.org", "", "http://localhost")
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
	fmt.Printf("%+v\n", message)*/
	u := url.URL{
		Scheme: "wss",
		Host: "echo.websocket.org",
		Path: "/",
	}
	fmt.Println("HERE")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("error dialing websocket; %v", err)
	}
	fmt.Println("HERE")
	if err = c.WriteMessage(websocket.TextMessage, []byte("test123")); err != nil {
		log.Fatalf("error sending request; %v", err)
	}
	fmt.Println("HERE")
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatalf("error receiving response; %v", err)
	}
	fmt.Printf("%s\n", message)
}
