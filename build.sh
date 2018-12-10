#!/usr/bin/env sh

#go build -ldflags "-X github.com/icoco/ixkit-cli/core.Build=`git rev-parse HEAD`" -o ixkit main.go

go build -ldflags "" -o ixkit console.go main.go

