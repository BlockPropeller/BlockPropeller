.PHONY: setup tools mod generate fmt lint build-all test watch-api precommit integration graph

setup: tools
	mkdir .chainup
	git config core.hooksPath .githooks

tools:
	cd ~ && \
		go get github.com/mgechev/revive@v0.0.0-20190910172647-84deee41635a && \
		go get github.com/smartystreets/goconvey@v0.0.0-20190731233626-505e41936337 && \
		go get github.com/KyleBanks/depth/cmd/depth@v1.2.1 && \
		go get github.com/cespare/reflex@v0.2.0

mod:
	go mod tidy

generate:
	go generate ./...

fmt:
	go fmt ./...

lint:
	revive --config revive.toml -formatter friendly ./...

build-all:
	go build ./...

test:
	go test ./... -p 1 -short

watch-api:
	reflex -d fancy -c reflex.conf

precommit: mod generate fmt lint build-all test

integration:
	 go test ./integration
