package main

import (
	"context"
	"vivianlab/internal/app"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), app.Deploy); err != nil {
		panic(err)
	}
}
