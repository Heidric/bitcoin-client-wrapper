#!/bin/bash

cd ../src/
go get golang.org/x/sys/unix
GOARCH=amd64 GOOS=linux go build -o ../test/app/bin/app