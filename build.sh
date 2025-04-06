#!/bin/bash
#cleanup previous build
rm -f ./client
rm -f ./watcher

#build client app
go build -o ./client ./client_src/cmd/main.go

#build watcher app
go build -o ./watcher ./connection_watcher/cmd/main.go
