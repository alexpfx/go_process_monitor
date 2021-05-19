all: client server

client:
	go build -ldflags "-s -w" -o bin/client gpm-clients/cmd/main.go
server:
	go build -ldflags "-s -w" -o bin/server  gpm-server/cmd/main.go


