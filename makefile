default: run

# Run stage: Run the application with `go run`
run:
	@echo "Running the application..."
	go run main.go

# Build stage: Compile the Go project into an executable
setup:
	@echo "Building the application..."
	go mod download
