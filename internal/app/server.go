package app

import (
	"context"
	"sync"
	"vivianlab/internal/pkg/echo"
	"vivianlab/internal/pkg/login"
	"vivianlab/internal/pkg/utils"
	"vivianlab/web"

	"net/http"
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

	handler http.Handler
}

func Deploy(ctx context.Context, app *App) error {
	appHandler := http.NewServeMux()
	app.handler = appHandler
	app.rw_timeout = Timeout

	//app.Logger(ctx).Debug("vivian: [launch] mysql", "connection", app.database.Get().Init(ctx) == nil)
	app.Logger(ctx).Info("vivian: [launch] app", "address", app.listener)

	//TODO: change handle names to be discreet
	//TODO: handles called only by registerd users
	//TODO: ^addiitonal middleware to verify that the request is safe
	//TODO: test all curl commands to verify the data the user recieves
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)
	appHandler.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(web.WebUI))))
	appHandler.Handle("/login", weaver.InstrumentHandler("login", AccountLogin(ctx, app)))
	appHandler.Handle("/login/generatekey", GenerateTwoFactorAuth(ctx, app))
	appHandler.Handle("/login/verifykey", VerifyTwoFactorAuth(ctx, app))
	appHandler.Handle("/ws", HandleWebSocketTimestamp(ctx, app))

	//TODO: cant access webapp using ServeTLS on Multiprocess! Needs to be
	//reconfigured
	//return http.ServeTLS(app.listener, app.handler, "../../certificates/server.crt", "../../certificates/server.key")
	return http.Serve(app.listener, app.handler)
}
