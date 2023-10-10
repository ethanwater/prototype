package main

import (
	"context"
	"log"

	"github.com/ServiceWeaver/weaver"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubUserQuery interface {
	Query(context.Context) error
}

type githubQuery struct {
	weaver.Implements[GithubUserQuery]
}

//var (
//	GithubUserQueryRequests = metrics.NewCounter(
//		"add_count",
//		"number of github_user_query_requests called",
//	)
//)

func (gq *githubQuery) Query(ctx context.Context) error {
	var query struct {
		Viewer struct {
			Login     githubv4.String
			Email     githubv4.String
			AvatarURL githubv4.URI
		}
	}
	oauth2TokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)

	//httpClient := oauth2.NewClient(context.Background(), oauth2TokenSrc)
	graphqlClient := githubv4.NewClient(oauth2.NewClient(context.Background(), oauth2TokenSrc))

	err2 := graphqlClient.Query(context.Background(), &query, nil)
	if err2 != nil {
		log.Fatal(err2)
	}
	gq.Logger(ctx).Debug("GithubUserQuery", "query", query)
	//GithubUserQueryRequests.Inc()
	return nil
}
