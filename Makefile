all: client server

client:
	go build -o bin/client gpm-clients/cmd/main.go
server:
	go build -o bin/server  gpm-server/cmd/main.go


