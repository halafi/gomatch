#
# Makefile for log matching tool in golang
#
# After everything is done use "./jsonizer" to run the tool.
# Tested on Ubuntu 13.10.
#
all: getGolang build #clean

getGolang:
	sudo apt-get install golang-go
	export PATH=$PATH:/usr/local/go/bin
build:
	go build jsonizer.go
clean:
	sudo apt-get remove golang-go
