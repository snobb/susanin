TARGET := httprouter
CMDSRC := ./cmd/${TARGET}

all: fmt vet build

vet:
	go vet ./...

fmt:
	gofmt -l -w .

test:
	go test -cover ./...

build:
	go build -o ./bin/${TARGET} ${CMDSRC}

clean:
	go clean ./...
	-rm -rf bin

.PHONY: build clean vet
