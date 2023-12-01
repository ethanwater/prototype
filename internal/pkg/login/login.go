package login

import (
	"context"
	"net/http"
	"sync/atomic"

	"vivianlab/database"
	"vivianlab/internal/pkg/auth"
	"vivianlab/internal/pkg/cache"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

var (
	totalSuccessfulAccountLogins = metrics.NewCounter(
		"totalSuccesfulAccountLogins",
		"the total number of times impl.Login is called and succesfully retrieves an account",
	)
	totalFailedAccountLogins = metrics.NewCounter(
		"totalFailedAccountLogins",
		"the total number of times impl.Login is called and fails retrieving an account",
	)
)

type T interface {
	Login(context.Context, string, string) (bool, error)
	GenerateAuthKey2FA(context.Context) (string, error)
	VerifyAuthKey2FA(context.Context, string, string) (bool, error)
}

type impl struct {
	weaver.Implements[T]
	weaver.Unrouted
	cache weaver.Ref[cache.Cache]
	tfa   weaver.Ref[auth.T]
	db    weaver.Ref[database.Database]
}

type Account struct {
	weaver.AutoMarshal
	ID       int
	Alias    string
	Name     string
	Email    string
	Password string
	Tier     int
}

var activeAccountEmail atomic.Value

func (l *impl) Login(ctx context.Context, email string, password string) (bool, error) {
	log := l.Logger(ctx)

	if !auth.SanitizeEmailCheck(email) {
		log.Error("vivian: [error] invalid email address", "err", http.StatusBadRequest)
		return false, nil
	}

	fetchedAccount, err := l.db.Get().FetchAccount(ctx, email)
	if err != nil {
		log.Error("vivian: [error] failure fetching account, user does not exist", "err", http.StatusNotFound)
		return false, nil
	}

	activeAccountEmail.Store(fetchedAccount.Email)

	if resp, err := l.cache.Get().Get(ctx, email); err != nil {
		log.Error("vivian: [error] no cache found", "err", weaver.RemoteCallError)
	} else {
		if password == resp {
			totalSuccessfulAccountLogins.Inc()
			log.Debug("vivian: [ok] fetched account: ", "alias", fetchedAccount.Alias)
			return true, nil
		}
	}

	hashChannel := make(chan bool, 1)
	go func() {
		result := auth.VerfiyHashPassword(fetchedAccount.Password, password)
		hashChannel <- result
	}()

	if email == fetchedAccount.Email && <-hashChannel {
		totalSuccessfulAccountLogins.Inc()
		log.Debug("vivian: [ok] fetched account: ", "alias", fetchedAccount.Alias)
		if err := l.cache.Get().Put(ctx, email, password); err != nil {
			log.Error("vivian: [error] unable to cache", "err", weaver.RemoteCallError)
		}
		return true, nil
	} else {
		totalFailedAccountLogins.Inc()
		log.Error("vivian: [error] invalid credentials", "err", http.StatusBadRequest)
		return false, nil
	}
}

func (l *impl) GenerateAuthKey2FA(ctx context.Context) (string, error) {
	email := activeAccountEmail.Load().(string)
	authkey, err := l.tfa.Get().GenerateAuthKey2FA(ctx, email)

	return authkey, err
}

func (l *impl) VerifyAuthKey2FA(ctx context.Context, authkey_hash, input string) (bool, error) {
	result, err := l.tfa.Get().VerifyAuthKey2FA(ctx, authkey_hash, input)

	return result, err
}
