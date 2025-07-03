#!/bin/bash

cd "$(dirname $0)"

sed -n -i -re '0,/\/\/ ###/p' api.go
echo >> api.go
go run ./internal/gen | sed -re 's:(qerr|protocol)[.]:quic.:g' >> api.go
gofmt -w .
