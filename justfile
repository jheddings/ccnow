default:
    @just --list

# Install dependencies
install:
    npm install

# Build TypeScript
build:
    npm run build

# Run tests
test:
    npm test

# Run tests in watch mode
test-watch:
    npm test -- --watch

# Clean build artifacts
clean:
    rm -rf dist

# Build and run with sample input
dev:
    just build
    echo '{"cwd":"/tmp/test","context_window":{"used_percentage":42,"current_usage":{"input_tokens":38000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":1500}}}' | node dist/cli.js
