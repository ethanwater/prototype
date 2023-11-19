package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var calls atomic.Int32
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
				calls.Add(1)
				//app.Logger(ctx).Debug(fmt.Sprint(uint32(calls.Load())))
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

var liveConn *websocket.Conn
var socketSync sync.Mutex

func SocketCalls(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socketSync.Lock()
		if liveConn != nil { liveConn.Close() }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			app.Logger(ctx).Error("vivian: socket: [error] handshake failure", "err", websocket.HandshakeError{})
		} else {
			app.Logger(ctx).Info("vivian: socket: [ok] handshake success", "remote", conn.RemoteAddr(), "local", conn.LocalAddr())
		}
		liveConn = conn
		socketSync.Unlock()

		defer conn.Close()
			
		for {
			var once sync.Once
			select {
			case <-ctx.Done():
				app.Logger(ctx).Error("vivian: socket: [error]", "err", "context lost")
				return
			default:
				once.Do(func(){
					app.Logger(ctx).Debug("vivian: socket: [ok] timestamp.calls", "amt", uint32(calls.Load()))
				})
				err := liveConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(uint32(calls.Load()))))
				if err != nil {
					if err := app.utils.Get().LoggerSocket(ctx, "vivian: socket: [error] disconnected <- broken pipe ?"); err != nil {
						app.Logger(ctx).Error("vivian: socket: [error]", "err", err)
					}
					return
				}
				time.Sleep(time.Second)
			}
		}
	})
}
