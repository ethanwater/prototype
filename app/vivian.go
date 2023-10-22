package app

import (
	"context"
	"database/sql"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/pelletier/go-toml"
)

type App struct {
	weaver.Implements[weaver.Main]
	handles       weaver.Ref[handleInterface]
	listener      weaver.Listener `weaver:"vivian"`
	read_timeout  time.Duration
	write_timeout time.Duration

	handler  http.Handler
	database *sql.DB
}

func Deploy(ctx context.Context, app *App) error {
	toml, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("[vivian:%s]\n", toml.Get("vivian.version"))
	}

	appHandler := http.NewServeMux()
	app.handler = appHandler
	app.read_timeout = 10 * time.Second
	app.write_timeout = 10 * time.Second
	app.database = EstablishLinkDatabase(ctx, app)

	app.Logger(ctx).Info("vivian: app deployed", "address", app.listener)

	appHandler.Handle("/", serverControls{}.base(ctx, app))
	appHandler.Handle("/kill", serverControls{}.kill(ctx, app))
	appHandler.Handle("/add", serverControls{}.add(ctx, app))
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)
	appHandler.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		echo, err := app.handles.Get().InnerEcho(r.Context(), query)
		if err != nil {
			app.Logger(r.Context()).Error("error getting query results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(echo)
		if err != nil {
			app.Logger(r.Context()).Error("error marshaling search results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("error writing search results", "err", err)
		}
		app.Logger(ctx).Debug("echo", "echo", echo)
	})

	return http.Serve(app.listener, app.handler)
}
