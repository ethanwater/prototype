package main

import (
	"context"
	_ "embed"
	"net/http"
	"prototype/query"

	"github.com/ServiceWeaver/weaver"
)

// net/http: https://pkg.go.dev/net/http#pkg-overview
// serviceweaver: https://serviceweaver.dev/docs.html#components

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		panic(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	//subscriber weaver.Ref[Subscriber]
	query    weaver.Ref[query.GithubUserQuery]
	listener weaver.Listener `weaver:"proto"`
}

func run(ctx context.Context, a *app) error {
	a.Logger(ctx).Info("prototype run", "addr:", a.listener)

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path != "/" {
	//		http.NotFound(w, r)
	//		return
	//	}
	//	if _, err := fmt.Fprint(w); err != nil {
	//		a.Logger(r.Context()).Error("error writing index.html", "err", err)
	//	}
	//})

	http.HandleFunc("/github_user_query", func(w http.ResponseWriter, r *http.Request) {
		a.query.Get().Query(ctx)
	})

	//novuClient := novu.NewAPIClient("93e56d580ce3386f02eb1ca728c0c2f2", &novu.Config{})
	//stat, err3 := a.subscriber.Get().InitSubscriber(ctx, novuClient, query.Viewer.Email)
	//if err3 != nil {
	//	return err3
	//}

	return http.Serve(a.listener, nil)
}
