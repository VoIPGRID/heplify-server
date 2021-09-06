NAME?=heplify-server

PKGLIST=$(shell go list ./... | grep -Ev '/vendor|/metric|/config|/sipparser/internal')


all:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(NAME) cmd/heplify-server/*.go

debug:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(NAME) cmd/heplify-server/*.go

test:
	go vet $(PKGLIST)
	go test $(PKGLIST) -race

.PHONY: clean
clean:
	rm -fr $(NAME)
