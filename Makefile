GOROOT := ~/go/
GOPATH := ~/go/bin
INSTALLPATH := /usr/local/bin/jsonizer

#all: build clean
all: install_dependencies build clean
#all: install_dependencies build clean uninstall

install_dependencies:
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
	sudo rm /usr/local/bin/jsonizer
	sudo apt-get remove golang-go
