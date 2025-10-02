.PHONY: build
build:
	GOARCH=arm64 go build -tags drm,pi -o bin/clock
