.PHONY: all linux windows mac clean

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null)
VER_NODOT := $(subst .,,${VERSION})
DIST    := dist

all: linux windows mac

dist:
	mkdir -p $(DIST)

linux:
	GOOS=linux GOARCH=amd64 go build -o ./${DIST}/quill_${VER_NODOT}_linux_amd64
	GOOS=linux GOARCH=386 go build -o ./${DIST}/quill_${VER_NODOT}_linux_386

windows:
	GOOS=windows GOARCH=amd64 go build -o ./${DIST}/quill_${VER_NODOT}_windows_amd64.exe
	GOOS=windows GOARCH=386 go build -o ./${DIST}/quill_${VER_NODOT}_windows_386.exe

mac:
	GOOS=darwin GOARCH=amd64 go build -o ./${DIST}/quill_${VER_NODOT}_macOS_amd64

clean:
	rm -rf $(DIST)
