#!/usr/bin/env bash

go build -ldflags "-X main.Build=46e82ab" -o $HOME/$GOBIN/fstree
