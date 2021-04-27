all: client server

client:
	go build -o bin/client cmd/client/main.go
server:
	go build -o bin/server cmd/server/main.go


