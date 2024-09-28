BINARY_NAME := sqlxp
BINARY_PATH := $(PWD)/bin/$(BINARY_NAME)
SYMLINK_PATH := /usr/local/bin/$(BINARY_NAME)
CMD_DIR := $(PWD)/cmd

clean:
	rm -rf ./bin

build:
	go build -o $(BINARY_PATH) $(CMD_DIR)

# Will only work on unix-like systems and used locally for development
link: clean build
	sudo ln -sf $(BINARY_PATH) $(SYMLINK_PATH)

unlink:
	sudo rm -f $(SYMLINK_PATH)
