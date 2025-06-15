#!/usr/bin/bash

echo "Linting ..."
golangci-lint run --timeout 5m --issues-exit-code 1 -v