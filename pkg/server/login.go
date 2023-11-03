package server

import (
	"context"
	"vivian/pkg/auth/utils"

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
		log.Error("vivian: ERROR! invalid email address")
		return false, nil
	}

	//TODO: implement password sanitization (follows password requirement rules)
	//if !utils.SanitizePasswordCheck(password) {
	//	log.Error("vivian: ERROR! invalid password")
	//	return false, nil
	//}

	fetchedAccount, err := FetchAccount(ctx, email)
	if err != nil {
		log.Error("vivian: ERROR! failure fetching account, user does not exist", "err", err)
	}

	if email == fetchedAccount.Email && utils.VerfiyHashPassword(fetchedAccount.Password, password) {
		log.Debug("vivian: SUCCESS! fetched account: ", "alias", fetchedAccount.Alias)
		return true, nil
	} else {
		log.Debug("vivian: ERROR! invalid credentials")
		return false, nil
	}
}
