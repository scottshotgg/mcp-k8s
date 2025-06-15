#!/usr/bin/bash

echo "Building router"
CGO_ENABLED=0 $(which go) build -o "./bin/router" "./cmd/router"