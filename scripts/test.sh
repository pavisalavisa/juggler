#!/bin/zsh

set -o errexit
set -o nounset

echo ">> Running unit tests..."
go test ./...


