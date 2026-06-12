#!/usr/bin/env bash

apt-get -y install golang

make test-coverage
