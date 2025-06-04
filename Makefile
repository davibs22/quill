.PHONY: all linux windows mac

all: linux windows mac

VERSION=v010

dist:
	mkdir dist2

linux:
	GOOS=linux GOARCH=amd64 go build -o ./dist2/quill_${VERSION}_linux_amd64
	GOOS=linux GOARCH=386 go build -o ./dist2/quill_${VERSION}_linux_386

windows:
	GOOS=windows GOARCH=amd64 go build -o ./dist2/quill_${VERSION}_windows_amd64.exe
	GOOS=windows GOARCH=386 go build -o ./dist2/quill_${VERSION}_windows_386.exe

mac:
	GOOS=darwin GOARCH=amd64 go build -o ./dist2/quill_${VERSION}_macOS_amd64