package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocketTimestamp(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			app.Logger(ctx).Error("vivian: socket:", "upgrade fail", err)
		}
		defer conn.Close()

		for {
			// Sending a timestamp to the client every second
			err := conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(time.RFC3339)))
			if err != nil {
				app.Logger(ctx).Error("vivian: socket:", err)
				break
			}
			time.Sleep(time.Second)
		}
	})
}
