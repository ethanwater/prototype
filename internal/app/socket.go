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
		} else {
			app.utils.Get().LoggerSocket(ctx)
		}
		defer conn.Close()

		for {
			// Sending a timestamp to the client every second
			rn, _ := app.utils.Get().Time(ctx)
			err := conn.WriteMessage(websocket.TextMessage, rn)
			if err != nil {
				app.Logger(ctx).Error("vivian: socket:", err)
				break
			}
			time.Sleep(time.Second)
		}
	})
}
