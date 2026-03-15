# justfile for ccglow

module := "github.com/jheddings/ccglow"
binary := "ccglow"

# run setup on first invocation
default: setup

# setup the local development environment
setup:
	go mod tidy

# build the binary
build:
	go build -o dist/{{binary}} .

# run tests
test:
	go test ./...

# auto-format
tidy: setup
	gofmt -w .
	npx --yes prettier --write .

# run format and vet checks
check:
	gofmt -l . | grep . && exit 1 || true
	npx --yes prettier --check .
	go vet ./...

# full preflight: build + check + test
preflight: build check test

# preview release notes since the last tag (or for a given range)
notes tag="--unreleased":
	npx git-cliff {{tag}}

# bump version, preflight, commit, tag, and push
release bump="patch": preflight
	#!/usr/bin/env bash
	# read current version from git tags
	CURRENT=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
	if [ -z "$CURRENT" ]; then
		CURRENT="0.0.0"
	fi
	IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT"
	case "{{bump}}" in
		major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
		minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
		patch) PATCH=$((PATCH + 1)) ;;
		*) echo "Unknown bump type: {{bump}}"; exit 1 ;;
	esac
	VERSION="$MAJOR.$MINOR.$PATCH"
	git commit --allow-empty -m "ccglow-$VERSION"
	git tag -a "v$VERSION" -m "v$VERSION"
	git push && git push --tags

# remove build artifacts
clean:
	rm -f {{binary}}

# remove everything including caches
clobber: clean
	go clean -cache
