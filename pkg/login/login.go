package login

import (
	"context"
	"net/http"
	"vivianlab/pkg/database"
	"vivianlab/pkg/utils"

	"github.com/ServiceWeaver/weaver"
)

type Login interface {
	Login(context.Context, string, string) (bool, error)
}

type login struct {
	weaver.Implements[Login]
}

func (l *login) Login(ctx context.Context, email string, password string) (bool, error) {
	log := l.Logger(ctx)

	if !utils.SanitizeEmailCheck(email) {
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
	if email == fetchedAccount.Email && utils.VerfiyHashPassword(fetchedAccount.Password, password) {
		log.Debug("vivian: SUCCESS! fetched account: ", "alias", fetchedAccount.Alias)
		return true, nil
	} else {
		log.Error("vivian: ERROR! invalid credentials", "err", http.StatusBadRequest)
		return false, nil
	}
}
