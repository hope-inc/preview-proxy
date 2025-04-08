GIT_VER := $(shell git describe --tags)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)

preview-proxy: cmd/*.go go.mod go.sum
	CGO_ENABLED=0 go build -o preview-proxy ./cmd/main.go

clean:
	rm -rf dist/* preview-proxy

run: preview-proxy
	./preview-proxy

packages:
	goreleaser release --rm-dist --snapshot --skip-publish

docker-image:
	docker build -t ghcr.io/hope-inc/preview-proxy:$(GIT_VER) -f Dockerfile .

push-image: docker-image
	docker push ghcr.io/hope-inc/preview-proxy:$(GIT_VER)

test:
	go test -v ./...

install:
	go install github.com/hope-inc/preview-proxy/cmd/preview-proxy
