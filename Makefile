#
# Multi-pattern searching in a single logfile.
#
# After everything is done use "jsonizer" to run the tool.
# Tested on Ubuntu 13.10.

GOROOT := ~/go/
GOPATH := ~/go/bin
INSTALLPATH := /usr/local/bin/jsonizer

all: build clean

dependencies:
	sudo apt-get install golang-go
	sudo apt-get install bzr
build:
	GOPATH=$(GOPATH) go get code.google.com/p/go.crypto/ssh/terminal
	GOPATH=$(GOPATH) go build jsonizer.go
	sudo cp ./jsonizer $(INSTALLPATH)
clean:
	rm -r -f $(GOROOT)
	rm ./jsonizer
uninstall:
	sudo rm ./jsonizer /usr/local/bin/jsonizer
