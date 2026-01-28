BINARY_NAME=taraDrop
WINDOWS_CC=x86_64-w64-mingw32-gcc

linux:
	go build -o $(BINARY_NAME) main.go

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WINDOWS_CC) go build -ldflags "-H=windowsgui" -o $(BINARY_NAME).exe main.go

all: linux windows

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
