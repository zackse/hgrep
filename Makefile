all: build

build:
	go build

install:
	go install

updatedeps:
	go get -u
	go mod tidy

setup:
	go mod download
