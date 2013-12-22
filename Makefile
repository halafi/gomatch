#
# Makefile for log matching tool in golang
#
# After everything is done use "./jsonizer" to run the tool.
# Tested on Ubuntu 13.10.
#

GOPATH := ~/go/bin

all: build clean

dependencies:
	sudo apt-get install golang-go
	sudo apt-get install bzr
build:
	GOPATH=$(GOPATH) go get code.google.com/p/go.crypto/ssh/terminal
	GOPATH=$(GOPATH) go build jsonizer.go
clean:
	rm -r -f $(GOPATH)
