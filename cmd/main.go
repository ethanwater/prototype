package main

import (
	"context"
	"vivian/pkg/app"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), app.Deploy); err != nil {
		panic(err)
	}
}
