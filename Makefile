#
# Makefile for log matching tool in golang
#
# After everything is done use "./jsonizer" to run the tool.
# Tested on Ubuntu 13.10.
#
GOPATH := ~/go/bin
GOROOT := ~/go

all: get_dependencies build

get_dependencies:
	sudo apt-get install golang-go
	sudo apt-get install bzr

build:
	GOPATH=$(GOPATH) go get labix.org/v2/pipe
	GOPATH=$(GOPATH) go build jsonizer.go
