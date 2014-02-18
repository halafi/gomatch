DEST := /usr/local/bin
TARGETNAME := gomatch

all: dependencies build

dependencies:
	sudo apt-get install golang-go

build:
	go build -o $(TARGETNAME) src/*

install:
	sudo cp ./$(TARGETNAME) $(DEST)/$(TARGETNAME)
	rm $(TARGETNAME)

uninstall:
	sudo rm -f $(DEST)/$(TARGETNAME)
