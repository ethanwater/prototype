package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		panic(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	//subscriber weaver.Ref[Subscriber]
	query weaver.Ref[GithubUserQuery]
}

func run(ctx context.Context, a *app) error {
	a.query.Get().Query(ctx)
	//novuClient := novu.NewAPIClient("93e56d580ce3386f02eb1ca728c0c2f2", &novu.Config{})
	//stat, err3 := a.subscriber.Get().InitSubscriber(ctx, novuClient, query.Viewer.Email)
	//if err3 != nil {
	//	return err3
	//}

	return nil
}
