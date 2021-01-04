.PHONY: build
build:
	go build -v ./cmd/tgsubscriber/main.go


.DEFAULT: build