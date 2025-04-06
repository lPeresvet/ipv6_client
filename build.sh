#!/bin/bash
#cleanup previous build
rm -f ./client

#build client app
go build -o ./client ./client_src/cmd/main.go
