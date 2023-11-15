package app

import (
	"context"
	"database/sql"
	"sync"
	"vivianlab/database"
	"vivianlab/internal/pkg/echo"
	"vivianlab/internal/pkg/login"
	"vivianlab/internal/pkg/utils"
	"vivianlab/web"

	"net/http"
	"os"
	"time"

	"github.com/ServiceWeaver/weaver"
)

const Timeout = 10 * time.Second

type App struct {
	weaver.Implements[weaver.Main]
	listener weaver.Listener `weaver:"vivian"`
	echo     weaver.Ref[echo.Echoer]
	login    weaver.Ref[login.Login]
	utils    weaver.Ref[utils.T]

	rw_timeout time.Duration
	mux        sync.Mutex

	handler  http.Handler
	database *sql.DB
}

func Deploy(ctx context.Context, app *App) error {
	appHandler := http.NewServeMux()
	app.handler = appHandler
	app.rw_timeout = Timeout

	//TODO: non repetitive database integration
	db, err := database.EstablishLinkDatabase(ctx)

	if err != nil {
		app.Logger(ctx).Error("vivian: ERROR! cannot connect to database", "status", http.StatusRequestTimeout)
		os.Exit(13)
	} else {
		app.database = db
		app.Logger(ctx).Debug("vivian: CONNECTED to database", "status", http.StatusOK)
	}
	//

	//function to verify all statuses are valid
	app.Logger(ctx).Info("vivian: APP DEPLOYED", "address", app.listener, "status", http.StatusOK)

	//TODO: change handle names to be discreet
	//TODO: handles called only by registerd users
	//TODO: addiiton middleware to verify that the request is safe
	//TODO: test all curl commands to verify the data the user recieves
	appHandler.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(web.WebUI))))
	//appHandler.Handle("/echo", EchoResponse(ctx, app))
	//appHandler.Handle("/fetch", FetchUsers(ctx, app))
	appHandler.Handle("/login", weaver.InstrumentHandler("login", AccountLogin(ctx, app)))
	appHandler.Handle("/login/generatekey", GenerateTwoFactorAuth(ctx, app))
	appHandler.Handle("/login/verifykey", VerifyTwoFactorAuth(ctx, app))
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)
	appHandler.Handle("/ws", HandleWebSocketTimestamp(ctx, app))
	
	//TODO: cant access webapp using ServeTLS on Multiprocess!
	return http.ServeTLS(app.listener, app.handler, "../../certificates/server.crt", "../../certificates/server.key")
	//return http.Serve(app.listener, app.handler)
}
