# ccglow - AI Development Context

Composable, spaceship-style statusline for Claude Code. Written in Go.

## Architecture

Segment tree model: atomic segments (single data point) and composite segments
(groups with `enabled` gating). Providers lazily fetch and cache external data.
The main pipeline: parse CLI → load render tree → read stdin →
resolve providers → depth-first render → styled output.

**Key concepts**:

- `SegmentNode` — configuration (what to render, how it looks)
- `Segment` — runtime behavior (how to produce a value)
- `Provider` — lazy, cached data fetcher (git, pwd, context)
- Presets — Go functions that return `[]SegmentNode` trees

## First Principles

- A segment renders a specific piece of data with style.

## Project Structure

- `main.go` — entry point, CLI (cobra), pipeline orchestration
- `internal/types/` — shared types (`SegmentNode`, `Style`)
- `internal/session/` — stdin JSON parsing
- `internal/config/` — JSON config file parsing
- `internal/segment/` — segment registry and implementations
- `internal/provider/` — provider registry and implementations
- `internal/render/` — tree traversal, provider resolution, output
- `internal/style/` — ANSI styling, color level control
- `internal/preset/` — named layouts (default, minimal, full)

## Development

**Key commands**:

- `go build ./...` — build
- `go test ./...` — run all tests
- `go vet ./...` — static analysis

**Adding a segment**: Add a case in `internal/segment/segment.go`'s registry.

**Adding a provider**: Create `internal/provider/<name>.go` implementing `Provider`,
register it in `internal/provider/provider.go`.

## Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <description>
```

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`

Scope is optional but encouraged (e.g. `fix(git): ...`, `feat(cli): ...`).

## Branch Naming

Use the same type prefixes as commits, followed by a short description:

```
<type>/<short-description>
```

Examples: `feat/color-themes`, `fix/token-formatting`, `chore/update-deps`

## Guardrails

**Do**:

- Follow TDD — write failing tests first, then implement
- Keep segments focused — one data point per atomic segment
- Return `""` from segments when there's no data (fail silent)
- Run `go vet ./... && go test ./...` before merging
- Work on feature branches

**Don't**:

- Push directly to main
- Never force-push to main
- Skip tests for new segments or providers
- Put styling logic in segments — segments return raw values, the renderer applies style
- Mutate global color state — use `style.SetColorLevel()` from `internal/style/`
