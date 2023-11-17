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
			app.Logger(ctx).Error("vivian: socket: [error] handshake failure", "err", "connection refused")
		} else {
			app.Logger(ctx).Info("vivian: socket: [ok] handshake success", "remote", conn.RemoteAddr(), "local", conn.LocalAddr())
		}
		defer conn.Close()

		reconnectChannel := make(chan int)
		defer close(reconnectChannel)

		go func() {
			for {
				select {
				case <-reconnectChannel:
					if err := app.utils.Get().LoggerSocket(ctx, "vivian: [ok] socket: reconnected"); err != nil {
						app.Logger(ctx).Error("vivian: socket: [error]", "err", err)
					}
					return
				case <-ctx.Done():
					app.Logger(ctx).Error("vivian: socket: [error]", "err", "context lost")
					return
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				app.Logger(ctx).Error("vivian: socket: [error]", "err", "context lost")
				return
			default:
				timestamp, _ := app.utils.Get().Time(ctx)
				err := conn.WriteMessage(websocket.TextMessage, timestamp)
				if err != nil {
					if err := app.utils.Get().LoggerSocket(ctx, "vivian: socket: [error] disconnected <- broken pipe ?"); err != nil {
						app.Logger(ctx).Error("vivian: socket: [error]", "err", err)
					}
					reconnectChannel <- 1
					return
				}
				time.Sleep(time.Second)
			}
		}
	})
}
