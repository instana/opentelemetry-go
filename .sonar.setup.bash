#!/usr/bin/env bash

apt-get -y install golang make

make test-coverage
