# ccglow - AI Development Context

Composable, spaceship-style statusline for Claude Code. Written in Go.

## Architecture

Segment tree model: atomic nodes (`expr` or `value`) and composite nodes
(groups with `enabled` gating and `when` conditions). Providers fetch external
data and return nested maps. The main pipeline: parse CLI → load render tree →
read stdin → resolve providers into env → depth-first render → styled output.

**Key concepts**:

- `SegmentNode` — configuration (what to render, how it looks)
- `expr` nodes — evaluate an expression against the provider env (e.g. `git.branch`)
- `value` nodes — render static text (literals, separators, newlines)
- `Provider` — data fetcher returning nested maps (e.g. `{"git": {"branch": "main"}}`)
- `eval.Eval()` — cached expr-lang evaluation for both `expr` fields and `when` conditions
- Presets — JSON files loaded via `embed.FS`

## First Principles

- A node renders a specific piece of data with style.

## Project Structure

- `main.go` — entry point, CLI (cobra), pipeline orchestration
- `internal/types/` — shared types (`SegmentNode`, `Style`, `ProviderResult`)
- `internal/session/` — stdin JSON parsing
- `internal/config/` — JSON config file parsing
- `internal/eval/` — expr-lang compilation, caching, evaluation
- `internal/provider/` — provider registry and implementations
- `internal/render/` — tree traversal, env building, output
- `internal/style/` — ANSI styling, color level control
- `internal/preset/` — named layouts (default, minimal, full, moonwalk, f1)

## Development

**Key commands**:

- `go build ./...` — build
- `go test ./...` — run all tests
- `go vet ./...` — static analysis

**Adding a provider**: Create `internal/provider/<name>.go` implementing `DataProvider`,
register it in `internal/provider/provider.go`. Return nested maps in `Values`.

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
- Keep nodes focused — one data point per atomic node
- Return zero values from providers when there's no data (fail silent)
- Run `go vet ./... && go test ./...` before merging
- Work on feature branches

**Don't**:

- Push directly to main
- Never force-push to main
- Skip tests for new providers
- Put styling logic in providers — providers return raw values, the renderer applies style
- Mutate global color state — use `style.SetColorLevel()` from `internal/style/`
