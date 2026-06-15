#!/usr/bin/env bash

apt-get -y install make curl

GO_VERSION=1.26.0
GO_ARCH=amd64

curl -LO https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
rm -rf /usr/local/go && \
tar -C /usr/local -xzf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
rm -f go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

export GOROOT="/usr/local/go"
export GOPATH="/root/go"
export PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

make test-coverage
