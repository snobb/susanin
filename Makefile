include git.mk

TARGET := susanin
EXAMPLESRC := ./examples

all: fmt vet test

vet:
	go vet ./...

fmt:
	gofmt -l -w .

test:
	go test -cover ./...

examples: vet fmt
	go build -o ./bin/${TARGET} ${EXAMPLESRC}

clean:
	go clean ./...
	-rm -rf bin

.PHONY: build clean vet test
