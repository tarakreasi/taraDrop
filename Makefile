BINARY_NAME=taraDrop
WINDOWS_CC=x86_64-w64-mingw32-gcc
LDFLAGS=-ldflags="-s -w"
WINDOWS_LDFLAGS=-ldflags="-H=windowsgui -s -w"

.PHONY: linux windows all clean release

linux:
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WINDOWS_CC) go build $(WINDOWS_LDFLAGS) -o $(BINARY_NAME).exe main.go

all: linux windows

release: all
	@echo "Binaries ready for release:"
	@ls -lh $(BINARY_NAME) $(BINARY_NAME).exe

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
