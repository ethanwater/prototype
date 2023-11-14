#!/bin/zsh

cd ../build/vivian-single/
go build .
export PATH="$PATH:$HOME/go/bin"
weaver single deploy singledeploy.toml
