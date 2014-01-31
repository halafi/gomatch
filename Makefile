all: build

get_dependencies:
	sudo apt-get install golang-go

build:
	go build -o jsonizer src/*
