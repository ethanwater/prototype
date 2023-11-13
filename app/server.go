package app

import (
	"context"
	"database/sql"
	"sync"
	"vivianlab/pkg/admincontrols"
	"vivianlab/pkg/database"
	"vivianlab/pkg/login"
	"vivianlab/web"

	"log"
	"os"
	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/pelletier/go-toml"
)

const Timeout = 10 * time.Second

type App struct {
	weaver.Implements[weaver.Main]
	listener weaver.Listener `weaver:"vivian"`
	echo     weaver.Ref[admincontrols.Echo]
	login    weaver.Ref[login.Login]

	rw_timeout time.Duration
	mu         sync.Mutex

	db_name  string
	handler  http.Handler
	database *sql.DB
}

func Deploy(ctx context.Context, app *App) error {
	toml, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	appHandler := http.NewServeMux()
	app.handler = appHandler
	app.rw_timeout = Timeout

	//TODO: non repetitive database integration
	app.db_name = toml.Get("database.name").(string)
	app.database, err = database.EstablishLinkDatabase(ctx)

	if err != nil {
		app.Logger(ctx).Error("vivian: ERROR! cannot connect to database", err)
		os.Exit(13)
	} else {
		app.Logger(ctx).Debug("vivian: CONNECTED to database", "database", app.db_name)
	}
	//

	//function to verify all statuses are valid
	app.Logger(ctx).Info("vivian: APP DEPLOYED", "address", app.listener)

	appHandler.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(web.WebUI))))
	appHandler.Handle("/echo", EchoResponse(ctx, app))
	appHandler.Handle("/fetch", FetchUsers(ctx, app))
	appHandler.Handle("/login", weaver.InstrumentHandler("login", AccountLogin(ctx, app)))
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)


	
	return http.Serve(app.listener, app.handler)
	//return http.ServeTLS(app.listener, app.handler, "certificates/server.crt", "certificates/server.key")
}
