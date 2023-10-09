package main

//
//import (
//	"context"
//	"log"
//
//	"github.com/ServiceWeaver/weaver"
//	"github.com/segmentio/ksuid"
//	"github.com/shurcooL/githubv4"
//
//	novu "github.com/novuhq/go-novu/lib"
//)
//
//type Subscriber interface {
//	InitSubscriber(context.Context, *novu.APIClient, githubv4.String) (novu.SubscriberResponse, error)
//}
//
//type subscriber struct {
//	weaver.Implements[Subscriber]
//	weaver.AutoMarshal
//}
//
//func (s *subscriber) InitSubscriber(ctx context.Context, novuClient *novu.APIClient, email githubv4.String) (novu.SubscriberResponse, error) {
//	subscriberID := ksuid.New()
//	subscriber := novu.SubscriberPayload{
//		Email: string(email),
//	}
//
//	resp, err := novuClient.SubscriberApi.Identify(ctx, subscriberID.String(), subscriber)
//	if err != nil {
//		log.Fatal("Subscriber error: ", err.Error())
//	}
//
//	return resp, nil
//}
//
