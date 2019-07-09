.PHONY: init tools mod generate fmt lint test

init: tools
	git config core.hooksPath .githooks

tools:
	cd ~ && \
		go get github.com/mgechev/revive@b70717f5395a29c099e82291e6fdf6168642faac && \
		go get github.com/smartystreets/goconvey@68dc04aab96ae4326137d6b77330c224063a927e

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

watch:
	goconvey .

precommit: mod generate fmt lint test
