.PHONY: all

all: build exec

build:
	go build -o ./bin/bandl.exe cmd/app/main.go

exec:
	cd ./bin; ./bandl.exe