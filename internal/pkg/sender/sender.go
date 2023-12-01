package sender

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/ServiceWeaver/weaver"
)

var sender = "vivian@vivian.com"

type T interface {
	SendVerificationCodeEmail(context.Context, string, string) error
}

type impl struct {
	weaver.Implements[T]
}

func (s *impl) SendVerificationCodeEmail(ctx context.Context, receiver string, verificationCode string) error {
	auth := smtp.PlainAuth("", sender, "", "smtp.gmail.com")

	msg := []byte(fmt.Sprintf("verification code: %s", verificationCode))
	err := smtp.SendMail("smtp.gmail.com:587", auth, sender, []string{receiver}, msg)
	if err != nil {
		s.Logger(ctx).Error("cannot send email", err)
	} else {
		s.Logger(ctx).Debug("verification code sent to you")
	}

	return err
}
