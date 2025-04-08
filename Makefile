.PHONY: test test-coverage clean

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with coverage and check for 100% coverage
test-coverage-check:
	go test -v -coverprofile=coverage.out ./...
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$coverage < 100" | bc -l) -eq 1 ]; then \
		echo "Test coverage is below 100% (current: $$coverage%)"; \
		exit 1; \
	fi
	echo "Test coverage is 100%"

# Clean up generated files
clean:
	rm -f coverage.out coverage.html 