PKG ?= ./...

checks:
	go vet $(PKG)

test:
	go test $(PKG)
