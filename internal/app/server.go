package app

import (
	"context"
	"sync"
	"vivianlab/internal/pkg/login"
	"vivianlab/internal/pkg/utils"
	"vivianlab/internal/pkg/utils/echo"
	"vivianlab/web"

	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type App struct {
	weaver.Implements[weaver.Main]
	listener weaver.Listener `weaver:"vivian"`
	echo     weaver.Ref[echo.T]
	login    weaver.Ref[login.T]
	utils    weaver.Ref[utils.T]

	rw_timeout time.Duration
	mux        sync.Mutex

	handler http.Handler
}

func Deploy(ctx context.Context, app *App) error {
	appHandler := http.NewServeMux()
	app.handler = appHandler
	app.rw_timeout = 10 * time.Second

	app.Logger(ctx).Info("vivian: [launch] app", "net", app.listener.Addr().Network(),
		"address", app.listener.Addr().String())

	//TODO: change handle names to be discreet
	//TODO: handles called only by registerd users
	//TODO: ^addiitonal middleware to verify that the request is safe
	//TODO: test all curl commands to verify the data the user recieves
	appHandler.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(web.WebUI))))
	appHandler.Handle("/login", weaver.InstrumentHandler("login", AccountLogin(ctx, app)))
	appHandler.Handle("/login/generatekey", GenerateTwoFactorAuth(ctx, app))
	appHandler.Handle("/login/verifykey", VerifyTwoFactorAuth(ctx, app))
	appHandler.Handle("/ws", HandleWebSocketTimestamp(ctx, app))
	appHandler.Handle("/wscalls", SocketCalls(ctx, app))
	appHandler.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)

	return http.ServeTLS(app.listener, app.handler, "../../certificates/server.crt", "../../certificates/server.key")
	//return http.Serve(app.listener, app.handler)
}
