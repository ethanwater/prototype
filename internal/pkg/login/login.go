package login

import (
	"context"
	"net/http"
	"vivianlab/database"
	"vivianlab/internal/pkg/auth"

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

type Login interface {
	Login(context.Context, string, string) (bool, error)
	GenerateAuthKey2FA(context.Context) (string, error)
	VerifyAuthKey2FA(context.Context, string, string) (bool, error)
}

type impl struct {
	weaver.Implements[Login]
	tfa weaver.Ref[auth.Authenticator]
}

func (l *impl) Login(ctx context.Context, email string, password string) (bool, error) {
	//TODO: count impl attempts. OK on JS end, but can still be curled via term
	log := l.Logger(ctx)

	if !auth.SanitizeEmailCheck(email) {
		log.Error("vivian: ERROR! invalid email address", "err", http.StatusBadRequest)
		return false, nil
	}

	//TODO: implement password sanitization (follows password requirement rules)
	//if !utils.SanitizePasswordCheck(password) {
	//	log.Error("vivian: ERROR! invalid password")
	//	return false, nil
	//}

	fetchedAccount, err := database.FetchAccount(ctx, email)
	if err != nil {
		log.Error("vivian: ERROR! failure fetching account, user does not exist", "err", http.StatusNotFound)
		return false, nil
	}

	//DAMN! this VerifyHash takes a while
	if email == fetchedAccount.Email && auth.VerfiyHashPassword(fetchedAccount.Password, password) {
		totalSuccessfulAccountLogins.Inc()
		log.Debug("vivian: SUCCESS! fetched account: ", "alias", fetchedAccount.Alias)
		return true, nil
	} else {
		totalFailedAccountLogins.Inc()
		log.Error("vivian: ERROR! invalid credentials", "err", http.StatusBadRequest)
		return false, nil
	}
}

func (l *impl) GenerateAuthKey2FA(ctx context.Context) (string, error) {
	authkey, err := l.tfa.Get().GenerateAuthKey2FA(ctx)

	return authkey, err
}

func (l *impl) VerifyAuthKey2FA(ctx context.Context, authkey_hash, input string) (bool, error) {
	result, err := l.tfa.Get().VerifyAuthKey2FA(ctx, authkey_hash, input)

	return result, err
}
