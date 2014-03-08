DEST := /usr/local/bin
TARGETNAME := gomatch
GOPATH := ~/gotmp

all: dependencies build clean

dependencies:
	sudo apt-get install golang-go

build:
	GOPATH=$(GOPATH) go get github.com/streadway/amqp
	GOPATH=$(GOPATH) go build -o $(TARGETNAME) src/*

clean:
	rm -r -f $(GOPATH)

install:
	sudo cp ./$(TARGETNAME) $(DEST)/$(TARGETNAME)
	rm $(TARGETNAME)

uninstall:
	sudo rm -f $(DEST)/$(TARGETNAME)
