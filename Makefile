DEST := /usr/local/bin
TARGETNAME := jsonizer

all: build

get_dependencies:
	sudo apt-get install golang-go

build:
	go build -o $(TARGETNAME) src/*

install:
	sudo cp ./$(TARGETNAME) $(DEST)/$(TARGETNAME)

uninstall:
	sudo rm -f $(DEST)/$(TARGETNAME)
