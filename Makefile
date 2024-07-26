.PHONY: all build run clean

all: build run-listener

build: 
	@echo "Building binary..."
	@go build -o bin/notifier-example

run: run-listener

run-listener:
	@echo "Running listener..."
	@./bin/notifier-example listen

run-emitter:
	@echo "Running emitter..."
	@./bin/notifier-example emmit

clean:
	@echo "Cleaning up..."
	@rm -rf bin