BINARY_NAME=myapp

all: build

build:
	go build -o $(BINARY_NAME) main.go

clean:
	go clean
	rm $(BINARY_NAME)

