.PHONY: run build clean

# Default target
run:
	go run main.go

# Build the example
build:
	go build -o firewall-example main.go

# Run the built binary
run-binary: build
	./firewall-example

# Clean up
clean:
	rm -f firewall-example

# Check if environment variables are set
check-env:
	@if [ -z "$$LATITUDE_API_TOKEN" ]; then \
		echo "Error: LATITUDE_API_TOKEN is not set"; \
		exit 1; \
	fi
	@if [ -z "$$LATITUDE_PROJECT" ]; then \
		echo "Error: LATITUDE_PROJECT is not set"; \
		exit 1; \
	fi
	@echo "Environment variables are correctly set"

# Run with environment check
run-safe: check-env run 