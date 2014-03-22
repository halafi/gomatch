@echo off
SET GOPATH=\tmplib

go get github.com/streadway/amqp

go build -o gomatch.exe .\src

rmdir \tmplib /s /q