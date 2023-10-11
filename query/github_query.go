package query

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubUserQuery interface {
	Query(context.Context, string) ([]string, error)
}

type GithubQuery struct {
	weaver.Implements[GithubUserQuery]
}
type (
	UserFragment struct {
		Bio   string
		Email string
	}
)

func (gq *GithubQuery) Query(ctx context.Context, user string) ([]string, error) {
	gq.Logger(ctx).Debug("GithubUserQuery", "query", user)
	x := map[string]interface{}{
		"user": githubv4.String(user),
	}

	var user_query struct {
		RepositoryOwner struct {
			Login        string
			UserFragment `graphql:"... on User"`
		} `graphql:"repositoryOwner(login: $user)"`
	}
	oauth2TokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)

	httpClient := oauth2.NewClient(context.Background(), oauth2TokenSrc)
	graphqlClient := githubv4.NewClient(httpClient)

	err2 := graphqlClient.Query(context.Background(), &user_query, x)
	if err2 != nil {
		gq.Logger(ctx).Error("GithubUserQuery couldn't fetch", "error:", err2)
	} else {
		gq.Logger(ctx).Debug("GithubUserQuery fetch", "result:", user_query.RepositoryOwner.Login)
	}

	if user_query.RepositoryOwner.Email == "" {
		gq.Logger(ctx).Error("GithubUserQuery couldn't fetch user email", "error: ", "null email value")
	}

	user_data := []string{
		user_query.RepositoryOwner.Login,
		user_query.RepositoryOwner.Bio,
		user_query.RepositoryOwner.Email, //email doesn't seem to always fetch...
	}

	return user_data, nil
}
