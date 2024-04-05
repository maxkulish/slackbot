ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJ_NAME=slackbot
BIN_DIR = $(ROOT_DIR)/bin
MAC_BIN_AMD = bin/$(PROJ_NAME)-darwin-amd64
MAC_BIN_ARM = bin/$(PROJ_NAME)-darwin-arm64
LINUX_BIN_AMD = bin/$(PROJ_NAME)-linux-amd64
LINUX_BIN_ARM = bin/$(PROJ_NAME)-linux-arm64

help: __help__

__help__:
	@echo make build - build go executables in the ./bin folder
	@echo make clean - delete executables, download project from github and build
	@echo make coverage - run test coverage and open html file with results
	@echo make benchmark - run benchmark tests with memory alocation
	@echo make test - run all tests in a project

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

build: rebuild

build-all: clean
	make build_mac_amd64
	make build_mac_arm64
	make build_linux_amd64
	make build_linux_arm64

rebuild:
	rm -rf $(BIN_DIR)/$(PROJ_NAME)
	go build -o $(BIN_DIR)/$(PROJ_NAME) main.go

build_mac_amd64:
	cd $(ROOT_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build --race -o $(MAC_BIN_AMD) cmd/main.go

build_mac_arm64:
	cd $(ROOT_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build --race -o $(MAC_BIN_ARM) cmd/main.go

build_linux_amd64:
	cd $(ROOT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(LINUX_BIN_AMD) cmd/main.go

build_linux_arm64:
	cd $(ROOT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(LINUX_BIN_ARM) cmd/main.go

clean:
	rm -rf ./bin

# ==============================================================================
# Running tests within the local computer

test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test-all: test lint vuln-check

test-all-race: test-race lint vuln-check

benchmark:
	cd $(ROOT_DIR)
	go test -bench . -benchmem

coverage:
	cd $(ROOT_DIR)
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

zip_linux:
	cd $(ROOT_DIR)
	zip $(PROJ_NAME)_linux.zip $(BIN_DIR)/linux/$(PROJ_NAME)

zip_mac:
	cd $(ROOT_DIR)
	zip $(PROJ_NAME)_mac.zip $(BIN_DIR)/macos/$(PROJ_NAME)

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all
