package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Split(r rune) bool {
	return r == '&' || r == ','
}

func GenerateTwoFactorAuth(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authkey, err := app.login.Get().GenerateAuthKey2FA(ctx)
		if err != nil {
			//https://go.dev/src/net/http/status.go
			app.Logger(ctx).Error("vivian: [error]", "failure generating authentication key", http.StatusBadRequest)
			//GenerateTwoFactorAuth(ctx, app) //may cause a race, TODO: check this line out
			return
		}

		bytes, err := json.Marshal(authkey)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "failure marshalling resullts", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "failure writing search results", err)
		}
	})
}

func VerifyTwoFactorAuth(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		hash := strings.TrimSpace(q.Get("hash"))
		input := strings.TrimSpace(q.Get("input"))

		result, err := app.login.Get().VerifyAuthKey2FA(ctx, hash, input)
		if err != nil {
			app.Logger(ctx).Error("vivian: [error]", "err", "unable to verify input")
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "err", err)
		}
	})
}

func EchoResponse(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		err := app.echo.Get().EchoResponse(ctx, query)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "err", http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(query)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: [error]", "err", err)
		}
	})
}

//func FetchUsers(ctx context.Context, app *App) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		result, err := database.FetchDatabaseData(ctx)
//		if err != nil {
//			app.Logger(r.Context()).Error("vivian: [error]", "err", "cannot add user")
//		}
//		bytes, err := json.Marshal(result)
//		if err != nil {
//			app.Logger(r.Context()).Error("vivian: [error] failure marshalling results", "err", err)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
//			app.Logger(r.Context()).Error("vivian: [error] failure writing search results", "err", err)
//		}
//
//	})
//}

func AccountLogin(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result bool
		q := r.URL.Query()
		email := strings.TrimSpace(q.Get("q"))
		password := strings.TrimSpace(q.Get("p"))

		app.mux.Lock()
		result, err := app.login.Get().Login(ctx, email, password)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: [error] failure logging", "err", http.StatusBadRequest)
		}
		app.mux.Unlock()

		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: [error] failure marshalling results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: [error] failure writing search results", "err", err)
		}

	})
}
