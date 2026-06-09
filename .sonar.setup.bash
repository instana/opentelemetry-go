#!/usr/bin/env bash

apt-get -y install golang

go test -v ./... -coverprofile=coverage.out
