package register

import (
	"context"
	"sync"
	"vivianlab/database"
	"vivianlab/internal/pkg/utils"

	"github.com/ServiceWeaver/weaver"
)

type Account struct {
	weaver.AutoMarshal
	ID       int
	Alias    string
	Name     string
	Email    string
	Password string
	Tier     int
}

type T interface {
	FetchRegistrationForm(context.Context, string, string, string, string, string) error
}

type Register struct {
	weaver.Implements[T]
	db weaver.Ref[database.Database]
}

var registerPool sync.Pool

func (r *Register) FetchRegistrationForm(ctx context.Context, name, email, password, alias string, tier int) error {
	accountID := utils.GenerateAccountID()
	newUser := Account{
		ID:       accountID,
		Alias:    alias,
		Name:     name,
		Email:    email,
		Password: password,
		Tier:     tier,
	}

	if err := r.DeliverAccountToDatabase(ctx, newUser); err != nil {
		return err
	}

	return nil
}

func (r *Register) DeliverAccountToDatabase(ctx context.Context, account Account) error {
	if err := r.db.Get().AddAccount(ctx, database.Account(account)); err != nil {
		return err
	}
	return nil
}
