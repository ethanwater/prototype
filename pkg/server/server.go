package server

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"vivian/pkg/frontend"

	"log"
	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/pelletier/go-toml"
)

const Timeout = 10 * time.Second

type Server struct {
	weaver.Implements[weaver.Main]
	echo     weaver.Ref[EchoInterface]
	add      weaver.Ref[AddUserInterface]
	listener weaver.Listener `weaver:"vivian"`

	read_timeout  time.Duration
	write_timeout time.Duration
	mu            sync.Mutex

	db_name  string
	handler  http.Handler
	Database *sql.DB
}

func Deploy(ctx context.Context, app *Server) error {
	toml, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	appHandler := http.NewServeMux()
	app.handler = appHandler
	app.read_timeout = Timeout
	app.write_timeout = Timeout
	go func() {
		app.db_name = toml.Get("database.name").(string)
		app.Database = EstablishLinkDatabase(ctx)
		app.Logger(ctx).Debug("connected to database: ", "database", app.db_name)
	}()

	app.Logger(ctx).Info("vivian: app deployed", "address", app.listener)

	appHandler.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(frontend.WebUI))))
	appHandler.Handle("/kill", kill(ctx, app))
	appHandler.Handle("/echo", Echo(ctx, app))
	appHandler.Handle("/add", Add(ctx, app))
	appHandler.Handle("/fetch", FetchUsers(ctx, app))
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)

	return http.Serve(app.listener, app.handler)
}

func kill(ctx context.Context, app *Server) http.Handler {
	app.mu.Lock()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger(ctx).Warn("vivian: kill server...", "status code", 0, http.ErrServerClosed)
		defer app.mu.Unlock()
		os.Exit(0)
	})
}
