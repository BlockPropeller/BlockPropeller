.PHONY: init tools mod generate fmt lint test

init: tools
	git config core.hooksPath .githooks

tools:
	cd ~ && \
		go get github.com/mgechev/revive@b70717f5395a29c099e82291e6fdf6168642faac

mod:
	go mod tidy

generate:
	go generate ./...

fmt:
	go fmt ./...

lint:
	revive --config revive.toml -formatter friendly ./...

test:
	go test ./... -p 1

precommit: mod generate fmt lint test
