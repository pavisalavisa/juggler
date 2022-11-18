#!/usr/bin/env bash

set -o errexit
set -o nounset

echo ">> Running unit tests..."
go test ./internal/...


