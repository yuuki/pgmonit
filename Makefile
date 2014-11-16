BIN = pgmonit

all: clean build test

build:
	go build -o $(BIN) github.com/y-uuki/pgmonit

test:
	go test ./...

clean:
	rm -f $(BIN)
	go clean

