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
			app.Logger(ctx).Info("vivian: socket: handshake success")
		}
		defer conn.Close()

		reconnectChannel := make(chan struct{})
		defer close(reconnectChannel)

		go func() {
			for {
				select {
				case <-reconnectChannel:
					app.utils.Get().LoggerSocket(ctx, "vivian: socket: reconnected")
					return
				case <-ctx.Done():
					app.Logger(ctx).Error("vivian: socket: disconnected", "err", "context lost")
					return
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				app.Logger(ctx).Error("vivian: socket: disconnected", "err", "context lost")
				return
			default:
				rn, _ := app.utils.Get().Time(ctx)
				err := conn.WriteMessage(websocket.TextMessage, rn)
				if err != nil {
					reconnectChannel <- struct{}{}
					app.utils.Get().LoggerSocket(ctx, "vivian: socket: disconnected <- broken pipe ?")
					return
				}
				time.Sleep(time.Second)
			}
		}
	})
}
