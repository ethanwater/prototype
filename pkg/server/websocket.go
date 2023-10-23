package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Socket(ctx context.Context, head http.Header) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		socket, err := upgrader.Upgrade(w, r, head)
		if err != nil {
			log.Panicln(err)
		}

		log.Println("connected client")
		err = socket.WriteMessage(1, []byte("hello user"))
		if err != nil {
			log.Println(err)
		}
		Reader(socket)
	})
}

func Reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		fmt.Println(p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
