BINARY_NAME=xkcd

all: build

build:
	go build -o $(BINARY_NAME) cmd/xkcd/main.go

clean:
	go clean
	rm $(BINARY_NAME)

