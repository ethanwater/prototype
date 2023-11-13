package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"vivianlab/pkg/database"
	"vivianlab/pkg/utils"
)

func Split(r rune) bool {
	return r == '&' || r == ','
}

func GenerateTwoFactorAuth(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authkey, err := utils.GenerateAuthKey2FA()
		if err != nil {
			app.Logger(ctx).Error("vivian: ERROR!", "err", err)
		}
		app.Logger(ctx).Debug(authkey)

		bytes, err := json.Marshal(authkey)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling resullts", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}
	})
}

func VerifyTwoFactorAuth(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		hash := strings.TrimSpace(q.Get("hash"))
		input := strings.TrimSpace(q.Get("input"))

		result := utils.VerifyAuthKey2FA(hash, input)
		if !result {
			app.Logger(ctx).Debug("false")
		} else {
			app.Logger(ctx).Debug("true")
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling resullts", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}
	})
}

func EchoResponse(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		err := app.echo.Get().EchoResponse(ctx, query)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR!", "err", http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(query)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling resullts", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}
	})
}

//func AddAccount(ctx context.Context, app *App) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		query := strings.TrimSpace(r.URL.Query().Get("q"))
//		result, err := app.add.Get().DatabaseAddAccount(ctx, query)
//		if err != nil {
//			app.Logger(r.Context()).Error("vivian: ERROR! cannot add user", "err", err)
//		}
//
//		bytes, err := json.Marshal(result)
//		if err != nil {
//			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling results", "err", err)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
//			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
//		}
//	})
//}

func FetchUsers(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := database.FetchDatabaseData(ctx)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR!", "err", "cannot add user")
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}

	})
}

func AccountLogin(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result bool
		q := r.URL.Query()
		email := strings.TrimSpace(q.Get("q"))
		password := strings.TrimSpace(q.Get("p"))

		app.mu.Lock()
		result, err := app.login.Get().Login(ctx, email, password)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure logging", "err", http.StatusBadRequest)
		}
		app.mu.Unlock()

		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}

	})
}
