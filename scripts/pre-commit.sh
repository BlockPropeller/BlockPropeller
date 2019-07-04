#!/bin/sh

set -e

echo "Running pre-commit tasks..."

go mod tidy
go generate ./...
go fmt ./...
go test ./...

echo "Pre-commit tasks complete!"
