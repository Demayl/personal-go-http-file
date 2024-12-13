all: test build run clean

BIN=bin/http_server

test:
	go test ./...

build:
	go build -o ${BIN} ./cmd/http_server
run:
	go build -o ${BIN} ./cmd/http_server
	echo ./${BIN} cmd/
	./${BIN} cmd/
clean:
	go clean
	rm ${BIN}
