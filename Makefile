BINARY_NAME=thanos
INSTALL_DIR=/usr/local/bin

.PHONY: build
build:
	go build -ldflags="-s -w -X main.version=1.0.0" -o thanos main.go
install: build
	mv $(BINARY_NAME) $(INSTALL_DIR)
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)