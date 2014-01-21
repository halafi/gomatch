GOROOT := ~/go
GOPATH := ~/go/bin

all: build clean

get_dependencies:
	sudo apt-get install golang-go
	sudo apt-get install bzr
build:
	GOPATH=$(GOPATH) go get code.google.com/p/go.crypto/ssh/terminal
	GOPATH=$(GOPATH) go build jsonizer.go
clean:
	rm -r -f $(GOROOT)

uninstall:
	rm ./jsonizer
