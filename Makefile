.PHONY: generate test bench clean all perf scrub update help

# Default target
help:
	@echo "Available targets:"
	@echo "  make generate  - Generate optimized Go code from data/*.json files"
	@echo "  make test      - Run all tests"
	@echo "  make bench     - Run benchmarks with memory stats"
	@echo "  make clean     - Remove all generated files"
	@echo "  make all       - Generate code and run tests"
	@echo "  make perf      - Generate code and run benchmarks"
	@echo "  make scrub     - Update country data from Wikipedia/UN (requires network)"
	@echo "  make update    - Fetch new data and regenerate code"

# Generate the optimized country code functions from data/*.json files
generate:
	go generate

# Run tests
test:
	go test -v

# Run benchmarks with memory stats
bench:
	go test -bench=. -benchmem

# Clean generated files
clean:
	rm -f *_generated.go

# Regenerate and test
all: generate test

# Check performance
perf: generate bench

# Fetch latest country data from Wikipedia and UN sources (manual update only)
# Note: This runs in a separate module to keep main go.mod clean
scrub:
	@echo "Fetching latest country data from Wikipedia and UN..."
	cd cmd/scrubber && go run main.go

# Fetch data and regenerate code
update: scrub generate
	@echo "Data updated and code regenerated successfully!"