.PHONY: build
build:
	@echo "Building ai-commit..."
	go build -ldflags "-s -w" -o builds/ai-commit cmd/ai-commit/main.go

.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o builds/ai-commit-linux-amd64 cmd/ai-commit/main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o builds/ai-commit-darwin-amd64 cmd/ai-commit/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o builds/ai-commit-darwin-arm64 cmd/ai-commit/main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o builds/ai-commit-windows-amd64.exe cmd/ai-commit/main.go

.PHONY: install
install: build
	@echo "Installing ai-commit..."
	cp builds/ai-commit $(GOPATH)/bin/ai-commit

.PHONY: clean
clean:
	@echo "Cleaning builds..."
	rm -rf builds/

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.DEFAULT_GOAL := build