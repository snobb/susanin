TARGET := susanin
EXAMPLESRC := ./examples

all: fmt vet build

vet:
	go vet ./...

fmt:
	gofmt -l -w .

test:
	go test -cover ./...

build:
	go install ./pkg/${TARGET}

examples: vet fmt
	go build -o ./bin/${TARGET} ${EXAMPLESRC}

clean:
	go clean ./...
	-rm -rf bin

.PHONY: build clean vet
