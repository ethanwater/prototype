#!/bin/zsh

cd ../build/vivian-multi/
go build .
export PATH="$PATH:$HOME/go/bin"
weaver multi deploy multideploy.toml
