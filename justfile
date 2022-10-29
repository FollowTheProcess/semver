COVERAGE_DATA := "coverage.out"
COVERAGE_HTML := "coverage.html"

# By default print the list of recipes
_default:
    @just --list

# Show all justfile variables
show:
    @just --evaluate

# Tidy up dependencies in go.mod and go.sum
tidy:
    go mod tidy

# Run go fmt on all project files
fmt:
    go fmt ./...

# Run all project unit tests
test *flags: fmt
    go test -race ./... {{ flags }}

# Lint the project and auto-fix errors if possible
lint: fmt
    golangci-lint run --fix

# Calculate test coverage and render the html
cover:
    go test -race -cover -covermode=atomic -coverprofile={{ COVERAGE_DATA }} ./...
    go tool cover -html={{ COVERAGE_DATA }} -o {{ COVERAGE_HTML }}
    open {{ COVERAGE_HTML }}

# Remove build artifacts and other project clutter
clean:
    go clean ./...
    rm -rf {{ COVERAGE_DATA }} {{ COVERAGE_HTML }}

# Run unit tests and linting in one go
check: test lint

# Print lines of code (for fun)
sloc:
    find . -name "*.go" | xargs wc -l | sort -nr
