package subscriber

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	novu "github.com/novuhq/go-novu/lib"
	"github.com/segmentio/ksuid"
)

type SubscriberInterface interface {
	InitSubscriber(context.Context, string) error
}
type Subscriber struct {
	weaver.Implements[SubscriberInterface]
}

func (s *Subscriber) InitSubscriber(ctx context.Context, email string) error {
	subscriberID := ksuid.New()
	novuClient := novu.NewAPIClient("93e56d580ce3386f02eb1ca728c0c2f2", &novu.Config{})

	subscriber := novu.SubscriberPayload{
		Email: string(email),
	}
	_, err := novuClient.SubscriberApi.Identify(ctx, subscriberID.String(), subscriber)
	if err != nil {
		s.Logger(ctx).Error("InitSubscriber could not subscribe user", "error:", err)
	} else {
		s.Logger(ctx).Debug("InitSubscriber user subscribed")
	}

	return nil
}
