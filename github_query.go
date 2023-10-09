package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ServiceWeaver/weaver"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var query struct {
	Viewer struct {
		Login     githubv4.String
		Email     githubv4.String
		AvatarURL githubv4.URI
	}
}

type GithubUserQuery interface {
	Query(context.Context) error
}

type githubQuery struct {
	weaver.Implements[GithubUserQuery]
}

func (gq *githubQuery) Query(ctx context.Context) error {
	oauth2TokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)

	httpClient := oauth2.NewClient(context.Background(), oauth2TokenSrc)
	graphqlClient := githubv4.NewClient(httpClient)

	err2 := graphqlClient.Query(context.Background(), &query, nil)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("    User:", query.Viewer.Login)
	fmt.Println("    Email:", query.Viewer.Email)
	fmt.Println("    Avatar:", query.Viewer.AvatarURL)

	return err2
}
