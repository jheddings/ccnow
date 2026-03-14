# justfile for ccnow

# run setup on first invocation
default: setup

# setup the local development environment
setup:
	npm install

# build TypeScript
build:
	npm run build

# run tests
test:
	npm test

# run tests in watch mode
test-watch:
	npm test -- --watch

# auto-format and lint-fix
tidy:
	npx prettier --write .
	npx eslint . --fix

# run format, lint, and type checks (no fix)
check:
	npx prettier --check .
	npx eslint .
	npx tsc --noEmit

# full preflight: build + check + test
preflight: build check test

# preview release notes since the last tag (or for a given range)
notes tag="--unreleased":
	npx git-cliff {{tag}}

# bump version, preflight, commit, tag, and push
release bump="patch": preflight
	#!/usr/bin/env bash
	npm version {{bump}} --no-git-tag-version
	VERSION=$(node -p "require('./package.json').version")
	npx prettier --write package.json package-lock.json
	git add package.json package-lock.json
	git commit -m "ccnow-$VERSION"
	git tag -a "v$VERSION" -m "v$VERSION"
	git push && git push --tags

# remove build artifacts
clean:
	rm -rf dist

# remove everything including node_modules
clobber: clean
	rm -rf node_modules
