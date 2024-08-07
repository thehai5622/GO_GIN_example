# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main hello.go

# Run the application
run:
	@echo "Run..."
	@go run hello.go

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main