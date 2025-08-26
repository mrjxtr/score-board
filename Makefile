GO := go
APP := score-board
BIN_DIR := bin
MAIN_PKG := .

# Strip symbols; keep binaries small. Windows adds -H=windowsgui to hide console.
LDFLAGS := -s -w

LINUX_BIN := $(BIN_DIR)/$(APP)-linux-amd64
WINDOWS_BIN := $(BIN_DIR)/$(APP)-windows-amd64.exe

.PHONY: all build linux windows clean

all: build

build: linux windows

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

linux: $(BIN_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(LINUX_BIN) $(MAIN_PKG)

windows: $(BIN_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(LDFLAGS) -H=windowsgui" -o $(WINDOWS_BIN) $(MAIN_PKG)

clean:
	rm -rf $(BIN_DIR)


