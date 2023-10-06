package main

import (
	"context"
	"fmt"
	"log"

	novu "github.com/novuhq/go-novu/lib"
)

func initSubscriber(novuClient *novu.APIClient) {
	ctx := context.Background()
	subscriberID := "1"
	subscriber := novu.SubscriberPayload{
		LastName: "ethanwater",
		Email:    "watermonic@gmail.com",
		Avatar:   "https://randomuser.me/api/portraits/thumb/women/79.jpg",
		Data: map[string]interface{}{
			"location": map[string]interface{}{
				"city":     "New York City",
				"state":    "New York",
				"country":  "United States",
				"postcode": "10001",
			},
		},
	}

	resp, err := novuClient.SubscriberApi.Identify(ctx, subscriberID, subscriber)
	if err != nil {
		log.Fatal("Subscriber error: ", err.Error())
		return
	}

	fmt.Println(resp)

}

func trigger(novuClient *novu.APIClient) {
	subscriberID := "1"
	eventId := "prototype-test"

	ctx := context.Background()
	to := map[string]interface{}{
		"lastName":     "Nwosu",
		"firstName":    "John",
		"subscriberId": subscriberID,
		"email":        "watermonic@gmail.com",
	}
	payload := map[string]interface{}{
		"name": "Hello World",
		"organization": map[string]interface{}{
			"logo": "https://happycorp.com/logo.png",
		},
	}
	triggerResp, err := novuClient.EventApi.Trigger(ctx, eventId, novu.ITriggerPayloadOptions{
		To:      to,
		Payload: payload,
	})
	if err != nil {
		log.Fatal("Novu error", err.Error())
		return
	}

	fmt.Println(triggerResp)
}
func main() {
	novuClient := novu.NewAPIClient("1", &novu.Config{})
	//initSubscriber(novuClient)
	trigger(novuClient)

}
