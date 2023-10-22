package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/pelletier/go-toml"
)

type App struct {
	weaver.Implements[weaver.Main]
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

	appHandler.Handle("/", serverControls{}.base())
	appHandler.Handle("/kill", serverControls{}.kill(ctx, app))
	appHandler.Handle("/echo", serverControls{}.echo(ctx, app))
	appHandler.Handle("/add", serverControls{}.add(ctx, app))
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)

	return http.Serve(app.listener, app.handler)
}
