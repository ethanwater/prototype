package main

import (
	"context"
	"vivian/pkg/server"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), server.Deploy); err != nil {
		panic(err)
	}
}
