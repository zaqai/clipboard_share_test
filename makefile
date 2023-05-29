VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

# run: build
# 	./clipboard_share

linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o clipboard_share-linux-amd64 $(GOFLAGS) .
linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o clipboard_share-linux-arm64 $(GOFLAGS) .
darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o clipboard_share-darwin-arm64 $(GOFLAGS) .

darwin-amd64:

windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o clipboard_share-windows-arm64 $(GOFLAGS) .



all: linux-amd64 linux-arm64 darwin-arm64 darwin-amd64 windows-amd64
	
