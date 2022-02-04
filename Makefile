all: build

tidy:
	go mod tidy

install: tidy
	go mod download && go mod verify

build: install
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -installsuffix cgo -o service main.go
