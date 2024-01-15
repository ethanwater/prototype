package app

import (
	"context"
	"sync"
	"vivianlab/internal/pkg/login"
	"vivianlab/internal/pkg/socket"
	"vivianlab/web"

	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/gorilla/mux"
)

type App struct {
	weaver.Implements[weaver.Main]
	listener weaver.Listener `weaver:"vivian"`
	login    weaver.Ref[login.T]
	socket   weaver.Ref[socket.T]

	rw_timeout time.Duration
	mux        sync.Mutex
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	webui := http.FileServer(http.FS(web.WebUI))
	router.PathPrefix("/").Handler(http.StripPrefix("/", webui))

	router.Handle("/login", weaver.InstrumentHandler("login", AccountLogin(r.Context(), app)))
	router.Handle("/login/generatekey", GenerateTwoFactorAuth(r.Context(), app)).Methods("GET")
	router.Handle("/login/verifykey", VerifyTwoFactorAuth(r.Context(), app))
	router.Handle("/ws", HandleWebSocketTimestamp(r.Context(), app))
	router.Handle("/wscalls", SocketCalls(r.Context(), app))
	router.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)

	router.ServeHTTP(w, r)
}

func Deploy(ctx context.Context, app *App) error {
	app.rw_timeout = 10 * time.Second

	app.Logger(ctx).Info("vivian: [launch] app", "net", app.listener.Addr().Network(),
		"address", app.listener.Addr().String())

	handler := &App{}

	return http.Serve(app.listener, handler)
}
