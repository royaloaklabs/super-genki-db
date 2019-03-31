CC := go
PWD := $(shell pwd)
OUTPUT_DIR := $(PWD)/bin
RELEASE := super-genki-db
VERSION ?= 0.1.3
ARCH ?= amd64
PLATFORMS := windows linux darwin
OS = $(word 1, $@)
DB_FILE = jisho-main.db

test:
	@$(CC) test ./...

clean:
	@rm -rf $(OUTPUT_DIR) $(DB_FILE) 2&>/dev/null

pre:
	@mkdir -p $(OUTPUT_DIR)

release: clean pre $(PLATFORMS)

windows: POSTFIX = .exe

local:
	go run main.go

install:
	@go install

$(PLATFORMS):
	echo "building $(RELEASE) for $(OS)"
	@GOOS=$(OS) GOARCH=$(ARCH) $(CC) build -o $(OUTPUT_DIR)/$(RELEASE)-$(VERSION)-$(OS)-$(ARCH)$(POSTFIX)

all: test clean release

.PHONY: test $(PLATFORMS) release pre clean all

