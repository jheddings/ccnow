# ccnow Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a working `npx -y ccnow` CLI that renders a composable, spaceship-style statusline for Claude Code by reading session JSON from stdin.

**Architecture:** Segment tree with atomic/composite nodes. DataProviders lazily resolve external data (git, pwd, context). A runner orchestrates the pipeline: parse CLI → load render tree → read stdin → resolve providers → depth-first render → styled output. DSL factory functions author presets internally; JSON config for user customization.

**Tech Stack:** TypeScript, Node.js, chalk (styling), jest (testing), justfile (project lifecycle)

**Spec:** `docs/superpowers/specs/2026-03-14-ccnow-design.md`

---

## File Structure

```
ccnow/
  package.json            # Package manifest, bin entry, dependencies
  tsconfig.json           # TypeScript compiler config
  jest.config.ts          # Jest config for TS
  .justfile               # Project lifecycle (build, test, lint, clean)
  README.md               # (exists)
  src/
    cli.ts                # Entry point: read stdin, parse args, run pipeline
    types.ts              # All shared interfaces (SegmentNode, Segment, etc.)
    session.ts            # Parse stdin JSON into SessionData
    registry.ts           # Segment registry (Map<string, Segment>)
    providers.ts          # DataProvider registry, lazy resolution, caching
    render.ts             # Tree traversal: enabled checks, render, style application
    runner.ts             # Pipeline orchestrator: ties CLI → tree → providers → render → output
    style.ts              # Apply StyleAttrs to a raw value string via chalk
    composites.ts         # Map CLI segment flags to composite render trees
    segments/
      literal.ts          # Static text segment
      sep.ts              # Separator segment
      pwd.name.ts         # Directory basename
      pwd.path.ts         # Full path
      pwd.smart.ts        # Truncated path (p10k-style, single strategy)
      git.branch.ts       # Current git branch
      git.insertions.ts   # Lines added
      git.deletions.ts    # Lines removed
      context.tokens.ts   # Token count (human-formatted)
      context.percent.ts  # Context window usage %
      index.ts            # Register all built-in segments
    providers/
      git.ts              # Shell out to git CLI
      pwd.ts              # Derive path variants from cwd
      context.ts          # Derive token/percent from session JSON
      index.ts            # Register all built-in providers
    dsl/
      index.ts            # Factory functions: StatusLine, Git, Pwd, Context, etc.
    presets/
      default.ts          # Default layout (pwd | git | context)
      minimal.ts          # Minimal layout (pwd | git.branch)
      full.ts             # Full layout (all segments, verbose)
      index.ts            # Preset registry
    config.ts             # Load and parse JSON config file into SegmentNode[]
    cli-parser.ts         # Parse process.argv: ordered segment flags + value flags
  tests/
    fixtures/
      session-basic.json  # Minimal session JSON for testing
      session-full.json   # Full session JSON with all fields
      session-no-git.json # Session JSON for non-git directory
    types.test.ts         # Type guard / validation tests
    session.test.ts       # Session parsing tests
    style.test.ts         # Style application tests
    render.test.ts        # Render tree traversal tests
    providers.test.ts     # Provider resolution and caching tests
    registry.test.ts      # Segment registry tests
    cli-parser.test.ts    # CLI arg parsing tests
    config.test.ts        # JSON config loading tests
    runner.test.ts        # End-to-end pipeline tests
    segments/
      literal.test.ts
      sep.test.ts
      pwd.test.ts         # All pwd.* segments
      git.test.ts         # All git.* segments
      context.test.ts     # All context.* segments
    providers/
      git.test.ts
      pwd.test.ts
      context.test.ts
    dsl/
      index.test.ts       # DSL factory function tests
    presets/
      default.test.ts
      minimal.test.ts
```

---

## Chunk 1: Project Scaffolding and Core Types

### Task 1: Initialize project with package.json, tsconfig, jest, justfile

**Files:**
- Create: `package.json`
- Create: `tsconfig.json`
- Create: `jest.config.ts`
- Create: `justfile`
- Create: `.gitignore`

- [ ] **Step 1: Create package.json**

```json
{
  "name": "ccnow",
  "version": "0.1.0",
  "description": "Composable, spaceship-style statusline for Claude Code",
  "type": "module",
  "main": "dist/cli.js",
  "bin": {
    "ccnow": "dist/cli.js"
  },
  "scripts": {
    "build": "tsc",
    "test": "NODE_OPTIONS='--experimental-vm-modules' jest"
  },
  "keywords": ["claude", "claude-code", "statusline", "cli"],
  "license": "MIT",
  "engines": {
    "node": ">=18"
  },
  "devDependencies": {
    "@types/jest": "^29.5.0",
    "@types/node": "^20.0.0",
    "jest": "^29.7.0",
    "ts-jest": "^29.2.0",
    "typescript": "^5.5.0"
  },
  "dependencies": {
    "chalk": "^5.3.0"
  }
}
```

- [ ] **Step 2: Create tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "Node16",
    "moduleResolution": "Node16",
    "outDir": "dist",
    "rootDir": "src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "declaration": true,
    "sourceMap": true
  },
  "include": ["src"],
  "exclude": ["node_modules", "dist", "tests"]
}
```

- [ ] **Step 3: Create jest.config.ts**

```ts
import type { Config } from 'jest';

const config: Config = {
  preset: 'ts-jest/presets/default-esm',
  testEnvironment: 'node',
  roots: ['<rootDir>/tests'],
  moduleNameMapper: {
    '^(\\.{1,2}/.*)\\.js$': '$1',
  },
  transform: {
    '^.+\\.tsx?$': [
      'ts-jest',
      {
        useESM: true,
      },
    ],
  },
};

export default config;
```

- [ ] **Step 4: Create justfile**

```makefile
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
```

- [ ] **Step 5: Create .gitignore**

```
node_modules/
dist/
*.tgz
.DS_Store
```

- [ ] **Step 6: Install dependencies**

Run: `cd /Users/jheddings/Projects/ccnow && npm install`
Expected: `node_modules` created, `package-lock.json` generated

- [ ] **Step 7: Verify TypeScript compiles (empty project)**

Run: `mkdir -p src && echo 'export {}' > src/cli.ts && npx tsc --noEmit`
Expected: No errors (remove placeholder after)

- [ ] **Step 8: Commit**

```bash
git add package.json package-lock.json tsconfig.json jest.config.ts justfile .gitignore
git commit -m "chore: scaffold project with package.json, tsconfig, jest, justfile"
```

---

### Task 2: Define core type interfaces

**Files:**
- Create: `src/types.ts`
- Create: `tests/fixtures/session-basic.json`
- Create: `tests/fixtures/session-full.json`
- Create: `tests/fixtures/session-no-git.json`

- [ ] **Step 1: Write type validation tests**

Create `tests/types.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import type { SegmentNode, Segment, SegmentContext, StyleAttrs, SessionData, DataProvider } from '../src/types.js';

describe('type contracts', () => {
  it('SegmentNode accepts atomic segment config', () => {
    const node: SegmentNode = {
      type: 'git.branch',
      provider: 'git',
      style: { color: 'white', bold: true },
    };
    expect(node.type).toBe('git.branch');
    expect(node.provider).toBe('git');
    expect(node.children).toBeUndefined();
  });

  it('SegmentNode accepts composite segment config', () => {
    const node: SegmentNode = {
      type: 'git',
      enabled: true,
      children: [
        { type: 'git.branch', provider: 'git' },
        { type: 'literal', props: { text: ' [' } },
      ],
    };
    expect(node.children).toHaveLength(2);
  });

  it('SegmentNode accepts enabled as function', () => {
    const node: SegmentNode = {
      type: 'git',
      enabled: (session) => session.cwd !== '',
      children: [],
    };
    expect(typeof node.enabled).toBe('function');
  });

  it('StyleAttrs supports all style properties', () => {
    const style: StyleAttrs = {
      color: 'cyan',
      bold: true,
      dim: false,
      italic: true,
      icon: '\ue0a0',
      prefix: '+',
      suffix: '%',
    };
    expect(style.color).toBe('cyan');
  });

  it('SessionData matches Claude Code stdin shape', () => {
    const session: SessionData = {
      cwd: '/Users/test/project',
      context_window: {
        used_percentage: 42,
        current_usage: {
          input_tokens: 38000,
          cache_creation_input_tokens: 2000,
          cache_read_input_tokens: 1500,
        },
      },
    };
    expect(session.cwd).toBe('/Users/test/project');
    expect(session.context_window?.used_percentage).toBe(42);
  });

  it('SessionData allows missing context_window', () => {
    const session: SessionData = { cwd: '/tmp' };
    expect(session.context_window).toBeUndefined();
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/types.test.ts`
Expected: FAIL — cannot find module `../src/types.js`

- [ ] **Step 3: Create src/types.ts**

```ts
export interface SessionData {
  cwd: string;
  context_window?: {
    used_percentage: number;
    current_usage?: {
      input_tokens: number;
      cache_creation_input_tokens: number;
      cache_read_input_tokens: number;
    };
  };
  [key: string]: unknown;
}

export interface StyleAttrs {
  color?: string;
  bold?: boolean;
  dim?: boolean;
  italic?: boolean;
  icon?: string;
  prefix?: string;
  suffix?: string;
}

export type EnabledFn = (session: SessionData) => boolean;

export interface SegmentNode {
  type: string;
  provider?: string;
  enabled?: boolean | EnabledFn;
  style?: StyleAttrs;
  props?: Record<string, unknown>;
  children?: SegmentNode[];
}

export interface SegmentContext {
  session: SessionData;
  provider?: unknown;
  props?: Record<string, unknown>;
}

export interface Segment {
  name: string;
  provider?: string;
  render(context: SegmentContext): string | null;
}

export interface DataProvider {
  name: string;
  resolve(session: SessionData): Promise<unknown>;
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/types.test.ts`
Expected: PASS — all 6 tests pass

- [ ] **Step 5: Create test fixtures**

Create `tests/fixtures/session-basic.json`:

```json
{
  "cwd": "/Users/test/project"
}
```

Create `tests/fixtures/session-full.json`:

```json
{
  "cwd": "/Users/test/project",
  "context_window": {
    "used_percentage": 42,
    "current_usage": {
      "input_tokens": 38000,
      "cache_creation_input_tokens": 2000,
      "cache_read_input_tokens": 1500
    }
  }
}
```

Create `tests/fixtures/session-no-git.json`:

```json
{
  "cwd": "/tmp",
  "context_window": {
    "used_percentage": 5,
    "current_usage": {
      "input_tokens": 5000,
      "cache_creation_input_tokens": 0,
      "cache_read_input_tokens": 0
    }
  }
}
```

- [ ] **Step 6: Commit**

```bash
git add src/types.ts tests/types.test.ts tests/fixtures/
git commit -m "feat: define core type interfaces and test fixtures"
```

---

### Task 3: Session parser

**Files:**
- Create: `src/session.ts`
- Create: `tests/session.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/session.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { parseSession } from '../src/session.js';

describe('parseSession', () => {
  it('parses valid JSON with all fields', () => {
    const input = JSON.stringify({
      cwd: '/Users/test/project',
      context_window: {
        used_percentage: 42,
        current_usage: { input_tokens: 38000, cache_creation_input_tokens: 2000, cache_read_input_tokens: 1500 },
      },
    });
    const session = parseSession(input);
    expect(session.cwd).toBe('/Users/test/project');
    expect(session.context_window?.used_percentage).toBe(42);
  });

  it('parses minimal JSON with only cwd', () => {
    const session = parseSession('{"cwd":"/tmp"}');
    expect(session.cwd).toBe('/tmp');
    expect(session.context_window).toBeUndefined();
  });

  it('returns null for empty string', () => {
    expect(parseSession('')).toBeNull();
  });

  it('returns null for invalid JSON', () => {
    expect(parseSession('not json')).toBeNull();
  });

  it('returns null for JSON missing cwd', () => {
    expect(parseSession('{"foo":"bar"}')).toBeNull();
  });

  it('preserves extra fields via index signature', () => {
    const session = parseSession('{"cwd":"/tmp","model":"opus"}');
    expect(session?.['model']).toBe('opus');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/session.test.ts`
Expected: FAIL — cannot find module `../src/session.js`

- [ ] **Step 3: Implement session parser**

Create `src/session.ts`:

```ts
import type { SessionData } from './types.js';

export function parseSession(input: string): SessionData | null {
  if (!input.trim()) return null;

  let parsed: unknown;
  try {
    parsed = JSON.parse(input);
  } catch {
    return null;
  }

  if (typeof parsed !== 'object' || parsed === null) return null;

  const obj = parsed as Record<string, unknown>;
  if (typeof obj.cwd !== 'string') return null;

  return obj as SessionData;
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/session.test.ts`
Expected: PASS — all 6 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/session.ts tests/session.test.ts
git commit -m "feat: add session JSON parser with validation"
```

---

### Task 4: Style application

**Files:**
- Create: `src/style.ts`
- Create: `tests/style.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/style.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { applyStyle } from '../src/style.js';
import type { StyleAttrs } from '../src/types.js';

describe('applyStyle', () => {
  it('returns value unchanged when no style attrs', () => {
    expect(applyStyle('hello', {})).toBe('hello');
    expect(applyStyle('hello', undefined)).toBe('hello');
  });

  it('applies prefix before value', () => {
    const result = applyStyle('42', { prefix: '+' });
    expect(result).toContain('+');
    expect(result).toContain('42');
    expect(result.indexOf('+')).toBeLessThan(result.indexOf('42'));
  });

  it('applies suffix after value', () => {
    const result = applyStyle('42', { suffix: '%' });
    expect(result).toContain('42');
    expect(result).toContain('%');
  });

  it('applies icon before prefix and value', () => {
    const result = applyStyle('main', { icon: '\ue0a0 ', prefix: '' });
    expect(result).toContain('\ue0a0');
    expect(result).toContain('main');
  });

  it('applies color via chalk (ANSI codes present)', () => {
    const result = applyStyle('hello', { color: 'cyan' });
    // chalk wraps with ANSI escape codes
    expect(result).toContain('hello');
    expect(result.length).toBeGreaterThan('hello'.length);
  });

  it('applies bold via chalk', () => {
    const result = applyStyle('hello', { bold: true });
    expect(result).toContain('hello');
    expect(result.length).toBeGreaterThan('hello'.length);
  });

  it('combines multiple style attrs', () => {
    const style: StyleAttrs = { color: 'green', bold: true, prefix: '+' };
    const result = applyStyle('12', style);
    expect(result).toContain('+');
    expect(result).toContain('12');
  });

  it('handles null/undefined style fields gracefully', () => {
    const style: StyleAttrs = { color: undefined, bold: undefined };
    expect(applyStyle('hello', style)).toBe('hello');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/style.test.ts`
Expected: FAIL — cannot find module `../src/style.js`

- [ ] **Step 3: Implement style application**

Create `src/style.ts`:

```ts
import chalk, { type ChalkInstance } from 'chalk';
import type { StyleAttrs } from './types.js';

export function applyStyle(value: string, style: StyleAttrs | undefined): string {
  if (!style) return value;

  // Build the decorated string: icon + prefix + value + suffix
  let result = value;
  if (style.prefix) result = style.prefix + result;
  if (style.suffix) result = result + style.suffix;
  if (style.icon) result = style.icon + result;

  // Apply chalk styling to the full string
  let painter: ChalkInstance = chalk;

  if (style.color) {
    // Support named colors and hex
    if (style.color.startsWith('#')) {
      painter = painter.hex(style.color);
    } else {
      painter = (painter as any)[style.color] ?? painter;
    }
  }
  if (style.bold) painter = painter.bold;
  if (style.dim) painter = painter.dim;
  if (style.italic) painter = painter.italic;

  // Only apply chalk if we actually set any style
  if (painter !== chalk) {
    result = painter(result);
  }

  return result;
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/style.test.ts`
Expected: PASS — all 8 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/style.ts tests/style.test.ts
git commit -m "feat: add style application with chalk integration"
```

---

## Chunk 2: Segment Registry, Providers, and Renderer

### Task 5: Segment registry

**Files:**
- Create: `src/registry.ts`
- Create: `tests/registry.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/registry.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { SegmentRegistry } from '../src/registry.js';
import type { Segment, SegmentContext } from '../src/types.js';

const mockSegment: Segment = {
  name: 'test.seg',
  render: (ctx: SegmentContext) => 'hello',
};

describe('SegmentRegistry', () => {
  it('registers and retrieves a segment by name', () => {
    const reg = new SegmentRegistry();
    reg.register(mockSegment);
    expect(reg.get('test.seg')).toBe(mockSegment);
  });

  it('returns undefined for unknown segment', () => {
    const reg = new SegmentRegistry();
    expect(reg.get('nope')).toBeUndefined();
  });

  it('registers multiple segments', () => {
    const reg = new SegmentRegistry();
    const seg2: Segment = { name: 'other', render: () => 'world' };
    reg.register(mockSegment);
    reg.register(seg2);
    expect(reg.get('test.seg')).toBe(mockSegment);
    expect(reg.get('other')).toBe(seg2);
  });

  it('later registration overwrites earlier', () => {
    const reg = new SegmentRegistry();
    const replacement: Segment = { name: 'test.seg', render: () => 'replaced' };
    reg.register(mockSegment);
    reg.register(replacement);
    expect(reg.get('test.seg')?.render({} as any)).toBe('replaced');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/registry.test.ts`
Expected: FAIL

- [ ] **Step 3: Implement segment registry**

Create `src/registry.ts`:

```ts
import type { Segment } from './types.js';

export class SegmentRegistry {
  private segments = new Map<string, Segment>();

  register(segment: Segment): void {
    this.segments.set(segment.name, segment);
  }

  get(name: string): Segment | undefined {
    return this.segments.get(name);
  }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/registry.test.ts`
Expected: PASS — all 4 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/registry.ts tests/registry.test.ts
git commit -m "feat: add segment registry"
```

---

### Task 6: Provider resolution and caching

**Files:**
- Create: `src/providers.ts`
- Create: `tests/providers.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/providers.test.ts`:

```ts
import { describe, it, expect, jest } from '@jest/globals';
import { ProviderRegistry } from '../src/providers.js';
import type { DataProvider, SessionData, SegmentNode } from '../src/types.js';

const session: SessionData = { cwd: '/tmp' };

const mockProvider: DataProvider = {
  name: 'test',
  resolve: async (s) => ({ value: s.cwd }),
};

describe('ProviderRegistry', () => {
  it('registers and resolves a provider', async () => {
    const reg = new ProviderRegistry();
    reg.register(mockProvider);
    const results = await reg.resolveAll(['test'], session);
    expect(results.get('test')).toEqual({ value: '/tmp' });
  });

  it('resolves multiple providers concurrently', async () => {
    const reg = new ProviderRegistry();
    const provider2: DataProvider = {
      name: 'other',
      resolve: async () => ({ data: 42 }),
    };
    reg.register(mockProvider);
    reg.register(provider2);
    const results = await reg.resolveAll(['test', 'other'], session);
    expect(results.get('test')).toEqual({ value: '/tmp' });
    expect(results.get('other')).toEqual({ data: 42 });
  });

  it('returns empty map when no providers requested', async () => {
    const reg = new ProviderRegistry();
    const results = await reg.resolveAll([], session);
    expect(results.size).toBe(0);
  });

  it('skips unknown provider names gracefully', async () => {
    const reg = new ProviderRegistry();
    const results = await reg.resolveAll(['nonexistent'], session);
    expect(results.has('nonexistent')).toBe(false);
  });

  it('catches provider errors and excludes from results', async () => {
    const reg = new ProviderRegistry();
    const failing: DataProvider = {
      name: 'fail',
      resolve: async () => { throw new Error('boom'); },
    };
    reg.register(failing);
    const results = await reg.resolveAll(['fail'], session);
    expect(results.has('fail')).toBe(false);
  });

  it('collects required provider names from a segment tree', () => {
    const reg = new ProviderRegistry();
    const tree: SegmentNode[] = [
      { type: 'pwd.smart', provider: 'pwd' },
      { type: 'sep' },
      { type: 'git', children: [
        { type: 'git.branch', provider: 'git' },
        { type: 'literal', props: { text: ' ' } },
        { type: 'git.insertions', provider: 'git' },
      ]},
    ];
    const names = reg.collectProviderNames(tree);
    expect([...names].sort()).toEqual(['git', 'pwd']);
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/providers.test.ts`
Expected: FAIL

- [ ] **Step 3: Implement provider registry**

Create `src/providers.ts`:

```ts
import type { DataProvider, SessionData, SegmentNode } from './types.js';

export class ProviderRegistry {
  private providers = new Map<string, DataProvider>();

  register(provider: DataProvider): void {
    this.providers.set(provider.name, provider);
  }

  collectProviderNames(tree: SegmentNode[]): Set<string> {
    const names = new Set<string>();
    const walk = (nodes: SegmentNode[]) => {
      for (const node of nodes) {
        // Skip statically disabled nodes (and their children)
        if (node.enabled === false) continue;
        if (node.provider) names.add(node.provider);
        if (node.children) walk(node.children);
      }
    };
    walk(tree);
    return names;
  }

  async resolveAll(
    names: string[],
    session: SessionData,
  ): Promise<Map<string, unknown>> {
    const results = new Map<string, unknown>();

    const entries = names
      .map((name) => ({ name, provider: this.providers.get(name) }))
      .filter((e): e is { name: string; provider: DataProvider } => e.provider !== undefined);

    const settled = await Promise.allSettled(
      entries.map(async ({ name, provider }) => {
        const data = await provider.resolve(session);
        return { name, data };
      }),
    );

    for (const result of settled) {
      if (result.status === 'fulfilled') {
        results.set(result.value.name, result.value.data);
      }
    }

    return results;
  }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/providers.test.ts`
Expected: PASS — all 6 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/providers.ts tests/providers.test.ts
git commit -m "feat: add provider registry with lazy resolution and caching"
```

---

### Task 7: Render tree traversal

**Files:**
- Create: `src/render.ts`
- Create: `tests/render.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/render.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { renderTree } from '../src/render.js';
import { SegmentRegistry } from '../src/registry.js';
import type { SessionData, SegmentNode, Segment } from '../src/types.js';

const session: SessionData = { cwd: '/Users/test/project' };
const providerData = new Map<string, unknown>();

function makeRegistry(...segments: Segment[]): SegmentRegistry {
  const reg = new SegmentRegistry();
  for (const seg of segments) reg.register(seg);
  return reg;
}

describe('renderTree', () => {
  it('renders a single atomic segment', () => {
    const reg = makeRegistry({ name: 'literal', render: (ctx) => (ctx.props?.text as string) ?? null });
    const tree: SegmentNode[] = [{ type: 'literal', props: { text: 'hello' } }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('hello');
  });

  it('applies style to segment output', () => {
    const reg = makeRegistry({ name: 'literal', render: (ctx) => (ctx.props?.text as string) ?? null });
    const tree: SegmentNode[] = [{ type: 'literal', props: { text: 'hello' }, style: { bold: true } }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toContain('hello');
    expect(result.length).toBeGreaterThan('hello'.length); // ANSI codes added
  });

  it('skips segments that render null', () => {
    const reg = makeRegistry({ name: 'empty', render: () => null });
    const tree: SegmentNode[] = [{ type: 'empty' }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('skips segments with enabled=false', () => {
    const reg = makeRegistry({ name: 'literal', render: (ctx) => (ctx.props?.text as string) ?? null });
    const tree: SegmentNode[] = [{ type: 'literal', props: { text: 'hidden' }, enabled: false }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('evaluates enabled function', () => {
    const reg = makeRegistry({ name: 'literal', render: (ctx) => (ctx.props?.text as string) ?? null });
    const tree: SegmentNode[] = [
      { type: 'literal', props: { text: 'shown' }, enabled: (s) => s.cwd === '/Users/test/project' },
      { type: 'literal', props: { text: 'hidden' }, enabled: (s) => s.cwd === '/nope' },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('shown');
  });

  it('renders composite children and concatenates', () => {
    const reg = makeRegistry({ name: 'literal', render: (ctx) => (ctx.props?.text as string) ?? null });
    const tree: SegmentNode[] = [{
      type: 'group',
      children: [
        { type: 'literal', props: { text: 'a' } },
        { type: 'literal', props: { text: 'b' } },
      ],
    }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('ab');
  });

  it('collapses composite to empty when all children are null', () => {
    const reg = makeRegistry({ name: 'empty', render: () => null });
    const tree: SegmentNode[] = [{
      type: 'group',
      children: [
        { type: 'empty' },
        { type: 'empty' },
      ],
    }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('collapses composite with disabled enabled flag', () => {
    const reg = makeRegistry({ name: 'literal', render: (ctx) => (ctx.props?.text as string) ?? null });
    const tree: SegmentNode[] = [{
      type: 'group',
      enabled: false,
      children: [
        { type: 'literal', props: { text: 'should not appear' } },
      ],
    }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('passes provider data to segment context', () => {
    const providers = new Map<string, unknown>([['git', { branch: 'main' }]]);
    const reg = makeRegistry({
      name: 'git.branch',
      provider: 'git',
      render: (ctx) => (ctx.provider as any)?.branch ?? null,
    });
    const tree: SegmentNode[] = [{ type: 'git.branch', provider: 'git' }];
    const result = renderTree(tree, reg, session, providers);
    expect(result).toBe('main');
  });

  it('skips segment when unknown type', () => {
    const reg = makeRegistry();
    const tree: SegmentNode[] = [{ type: 'unknown.thing' }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/render.test.ts`
Expected: FAIL

- [ ] **Step 3: Implement render tree traversal**

Create `src/render.ts`:

```ts
import type { SegmentNode, SessionData } from './types.js';
import { SegmentRegistry } from './registry.js';
import { applyStyle } from './style.js';

function isEnabled(node: SegmentNode, session: SessionData): boolean {
  if (node.enabled === undefined) return true;
  if (typeof node.enabled === 'boolean') return node.enabled;
  try {
    return node.enabled(session);
  } catch {
    return false;
  }
}

function renderNode(
  node: SegmentNode,
  registry: SegmentRegistry,
  session: SessionData,
  providerData: Map<string, unknown>,
): string | null {
  if (!isEnabled(node, session)) return null;

  // Composite node: render children
  if (node.children) {
    const parts: string[] = [];
    for (const child of node.children) {
      const rendered = renderNode(child, registry, session, providerData);
      if (rendered !== null) parts.push(rendered);
    }
    if (parts.length === 0) return null;
    const joined = parts.join('');
    return applyStyle(joined, node.style);
  }

  // Atomic node: look up segment and render
  const segment = registry.get(node.type);
  if (!segment) return null;

  const context = {
    session,
    provider: node.provider ? providerData.get(node.provider) : undefined,
    props: node.props,
  };

  const value = segment.render(context);
  if (value === null) return null;

  return applyStyle(value, node.style);
}

export function renderTree(
  tree: SegmentNode[],
  registry: SegmentRegistry,
  session: SessionData,
  providerData: Map<string, unknown>,
): string {
  const parts: string[] = [];
  for (const node of tree) {
    const rendered = renderNode(node, registry, session, providerData);
    if (rendered !== null) parts.push(rendered);
  }
  return parts.join('');
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/render.test.ts`
Expected: PASS — all 10 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/render.ts tests/render.test.ts
git commit -m "feat: add render tree traversal with enabled gating and style application"
```

---

## Chunk 3: Built-in Segments and Providers

### Task 8: Literal and separator segments

**Files:**
- Create: `src/segments/literal.ts`
- Create: `src/segments/sep.ts`
- Create: `tests/segments/literal.test.ts`
- Create: `tests/segments/sep.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/segments/literal.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { literalSegment } from '../../src/segments/literal.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('literal segment', () => {
  it('returns text from props', () => {
    expect(literalSegment.render({ session, props: { text: 'hello' } })).toBe('hello');
  });

  it('returns null when no text prop', () => {
    expect(literalSegment.render({ session, props: {} })).toBeNull();
  });

  it('returns null when no props', () => {
    expect(literalSegment.render({ session })).toBeNull();
  });

  it('has correct name', () => {
    expect(literalSegment.name).toBe('literal');
  });
});
```

Create `tests/segments/sep.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { sepSegment } from '../../src/segments/sep.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('sep segment', () => {
  it('returns char from props with spaces', () => {
    expect(sepSegment.render({ session, props: { char: '|' } })).toBe(' | ');
  });

  it('defaults to pipe when no char prop', () => {
    expect(sepSegment.render({ session })).toBe(' | ');
  });

  it('has correct name', () => {
    expect(sepSegment.name).toBe('sep');
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npm test -- tests/segments/`
Expected: FAIL

- [ ] **Step 3: Implement segments**

Create `src/segments/literal.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';

export const literalSegment: Segment = {
  name: 'literal',
  render(context: SegmentContext): string | null {
    const text = context.props?.text;
    if (typeof text !== 'string') return null;
    return text;
  },
};
```

Create `src/segments/sep.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';

export const sepSegment: Segment = {
  name: 'sep',
  render(context: SegmentContext): string | null {
    const char = (context.props?.char as string) ?? '|';
    return ` ${char} `;
  },
};
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npm test -- tests/segments/`
Expected: PASS — all 7 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/segments/literal.ts src/segments/sep.ts tests/segments/
git commit -m "feat: add literal and separator segments"
```

---

### Task 9: Pwd provider and segments

**Files:**
- Create: `src/providers/pwd.ts`
- Create: `src/segments/pwd.name.ts`
- Create: `src/segments/pwd.path.ts`
- Create: `src/segments/pwd.smart.ts`
- Create: `tests/providers/pwd.test.ts`
- Create: `tests/segments/pwd.test.ts`

- [ ] **Step 1: Write failing provider test**

Create `tests/providers/pwd.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { pwdProvider } from '../../src/providers/pwd.js';
import type { SessionData } from '../../src/types.js';

describe('pwd provider', () => {
  it('resolves name, path, and smart from cwd', async () => {
    const session: SessionData = { cwd: '/Users/jheddings/Projects/ccnow' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.name).toBe('ccnow');
    expect(data.path).toBe('/Users/jheddings/Projects/ccnow');
    expect(data.smart).toBeDefined();
  });

  it('handles root path', async () => {
    const session: SessionData = { cwd: '/' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.name).toBe('/');
    expect(data.path).toBe('/');
  });

  it('smart truncates long paths', async () => {
    const session: SessionData = { cwd: '/Users/jheddings/Projects/very/deep/nested/path' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart.length).toBeLessThan(data.path.length);
    expect(data.smart).toContain('path'); // always keeps last component
  });

  it('smart keeps short paths as-is', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('/tmp');
  });

  it('replaces home dir with ~ in smart', async () => {
    const home = process.env.HOME ?? '/Users/test';
    const session: SessionData = { cwd: `${home}/Projects/ccnow` };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart.startsWith('~')).toBe(true);
  });

  it('has correct name', () => {
    expect(pwdProvider.name).toBe('pwd');
  });
});
```

- [ ] **Step 2: Write failing segment tests**

Create `tests/segments/pwd.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { pwdNameSegment } from '../../src/segments/pwd.name.js';
import { pwdPathSegment } from '../../src/segments/pwd.path.js';
import { pwdSmartSegment } from '../../src/segments/pwd.smart.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/Users/test/project' };

describe('pwd.name segment', () => {
  it('returns directory name from provider data', () => {
    const provider = { name: 'project', path: '/Users/test/project', smart: '~/t/project' };
    expect(pwdNameSegment.render({ session, provider })).toBe('project');
  });

  it('returns null when no provider data', () => {
    expect(pwdNameSegment.render({ session })).toBeNull();
  });

  it('declares pwd provider', () => {
    expect(pwdNameSegment.provider).toBe('pwd');
  });
});

describe('pwd.path segment', () => {
  it('returns full path from provider data', () => {
    const provider = { name: 'project', path: '/Users/test/project', smart: '~/t/project' };
    expect(pwdPathSegment.render({ session, provider })).toBe('/Users/test/project');
  });
});

describe('pwd.smart segment', () => {
  it('returns smart-truncated path from provider data', () => {
    const provider = { name: 'project', path: '/Users/test/project', smart: '~/t/project' };
    expect(pwdSmartSegment.render({ session, provider })).toBe('~/t/project');
  });
});
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `npm test -- tests/providers/pwd.test.ts tests/segments/pwd.test.ts`
Expected: FAIL

- [ ] **Step 4: Implement pwd provider**

Create `src/providers/pwd.ts`:

```ts
import path from 'node:path';
import os from 'node:os';
import type { DataProvider, SessionData } from '../types.js';

export interface PwdData {
  name: string;
  path: string;
  smart: string;
}

function smartTruncate(cwd: string): string {
  const home = os.homedir();
  let p = cwd;

  // Replace home dir with ~
  if (p.startsWith(home)) {
    p = '~' + p.slice(home.length);
  }

  const parts = p.split('/');
  if (parts.length <= 3) return p;

  // Keep first component and last component, truncate middle to initials
  const first = parts[0]; // '' for absolute, '~' for home-relative
  const last = parts[parts.length - 1];
  const middle = parts.slice(1, -1).map((part) => part[0] ?? '');

  return [first, ...middle, last].join('/');
}

export const pwdProvider: DataProvider = {
  name: 'pwd',
  async resolve(session: SessionData): Promise<PwdData> {
    const cwd = session.cwd;
    return {
      name: cwd === '/' ? '/' : path.basename(cwd),
      path: cwd,
      smart: smartTruncate(cwd),
    };
  },
};
```

- [ ] **Step 5: Implement pwd segments**

Create `src/segments/pwd.name.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { PwdData } from '../providers/pwd.js';

export const pwdNameSegment: Segment = {
  name: 'pwd.name',
  provider: 'pwd',
  render(context: SegmentContext): string | null {
    const data = context.provider as PwdData | undefined;
    return data?.name ?? null;
  },
};
```

Create `src/segments/pwd.path.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { PwdData } from '../providers/pwd.js';

export const pwdPathSegment: Segment = {
  name: 'pwd.path',
  provider: 'pwd',
  render(context: SegmentContext): string | null {
    const data = context.provider as PwdData | undefined;
    return data?.path ?? null;
  },
};
```

Create `src/segments/pwd.smart.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { PwdData } from '../providers/pwd.js';

export const pwdSmartSegment: Segment = {
  name: 'pwd.smart',
  provider: 'pwd',
  render(context: SegmentContext): string | null {
    const data = context.provider as PwdData | undefined;
    return data?.smart ?? null;
  },
};
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `npm test -- tests/providers/pwd.test.ts tests/segments/pwd.test.ts`
Expected: PASS — all 11 tests pass

- [ ] **Step 7: Commit**

```bash
git add src/providers/pwd.ts src/segments/pwd.name.ts src/segments/pwd.path.ts src/segments/pwd.smart.ts tests/providers/pwd.test.ts tests/segments/pwd.test.ts
git commit -m "feat: add pwd provider and segments (name, path, smart)"
```

---

### Task 10: Context provider and segments

**Files:**
- Create: `src/providers/context.ts`
- Create: `src/segments/context.tokens.ts`
- Create: `src/segments/context.percent.ts`
- Create: `tests/providers/context.test.ts`
- Create: `tests/segments/context.test.ts`

- [ ] **Step 1: Write failing provider test**

Create `tests/providers/context.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { contextProvider } from '../../src/providers/context.js';
import type { SessionData } from '../../src/types.js';

describe('context provider', () => {
  it('formats tokens as K for thousands', async () => {
    const session: SessionData = {
      cwd: '/tmp',
      context_window: {
        used_percentage: 42,
        current_usage: { input_tokens: 38000, cache_creation_input_tokens: 2000, cache_read_input_tokens: 1500 },
      },
    };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBe('42K');
    expect(data.percent).toBe(42);
  });

  it('formats tokens as M for millions', async () => {
    const session: SessionData = {
      cwd: '/tmp',
      context_window: {
        used_percentage: 85,
        current_usage: { input_tokens: 1000000, cache_creation_input_tokens: 100000, cache_read_input_tokens: 100000 },
      },
    };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBe('1.2M');
  });

  it('returns raw number for small token counts', async () => {
    const session: SessionData = {
      cwd: '/tmp',
      context_window: {
        used_percentage: 1,
        current_usage: { input_tokens: 500, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 },
      },
    };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBe('500');
  });

  it('returns null fields when context_window missing', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBeNull();
    expect(data.percent).toBeNull();
  });

  it('has correct name', () => {
    expect(contextProvider.name).toBe('context');
  });
});
```

- [ ] **Step 2: Write failing segment tests**

Create `tests/segments/context.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { contextTokensSegment } from '../../src/segments/context.tokens.js';
import { contextPercentSegment } from '../../src/segments/context.percent.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('context.tokens segment', () => {
  it('returns formatted token string', () => {
    const provider = { tokens: '42K', percent: 42 };
    expect(contextTokensSegment.render({ session, provider })).toBe('42K');
  });

  it('returns null when tokens is null', () => {
    const provider = { tokens: null, percent: null };
    expect(contextTokensSegment.render({ session, provider })).toBeNull();
  });

  it('declares context provider', () => {
    expect(contextTokensSegment.provider).toBe('context');
  });
});

describe('context.percent segment', () => {
  it('returns percent as string', () => {
    const provider = { tokens: '42K', percent: 42 };
    expect(contextPercentSegment.render({ session, provider })).toBe('42%');
  });

  it('returns null when percent is null', () => {
    const provider = { tokens: null, percent: null };
    expect(contextPercentSegment.render({ session, provider })).toBeNull();
  });
});
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `npm test -- tests/providers/context.test.ts tests/segments/context.test.ts`
Expected: FAIL

- [ ] **Step 4: Implement context provider**

Create `src/providers/context.ts`:

```ts
import type { DataProvider, SessionData } from '../types.js';

export interface ContextData {
  tokens: string | null;
  percent: number | null;
}

function formatTokens(total: number): string {
  if (total >= 1_000_000) {
    return `${(total / 1_000_000).toFixed(1)}M`;
  }
  if (total >= 1_000) {
    return `${Math.round(total / 1_000)}K`;
  }
  return `${total}`;
}

export const contextProvider: DataProvider = {
  name: 'context',
  async resolve(session: SessionData): Promise<ContextData> {
    const cw = session.context_window;
    if (!cw) {
      return { tokens: null, percent: null };
    }

    const usage = cw.current_usage;
    let totalTokens = 0;
    if (usage) {
      totalTokens =
        (usage.input_tokens ?? 0) +
        (usage.cache_creation_input_tokens ?? 0) +
        (usage.cache_read_input_tokens ?? 0);
    }

    return {
      tokens: usage ? formatTokens(totalTokens) : null,
      percent: cw.used_percentage ?? null,
    };
  },
};
```

- [ ] **Step 5: Implement context segments**

Create `src/segments/context.tokens.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { ContextData } from '../providers/context.js';

export const contextTokensSegment: Segment = {
  name: 'context.tokens',
  provider: 'context',
  render(context: SegmentContext): string | null {
    const data = context.provider as ContextData | undefined;
    return data?.tokens ?? null;
  },
};
```

Create `src/segments/context.percent.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { ContextData } from '../providers/context.js';

export const contextPercentSegment: Segment = {
  name: 'context.percent',
  provider: 'context',
  render(context: SegmentContext): string | null {
    const data = context.provider as ContextData | undefined;
    if (data?.percent == null) return null;
    return `${data.percent}%`;
  },
};
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `npm test -- tests/providers/context.test.ts tests/segments/context.test.ts`
Expected: PASS — all 10 tests pass

- [ ] **Step 7: Commit**

```bash
git add src/providers/context.ts src/segments/context.tokens.ts src/segments/context.percent.ts tests/providers/context.test.ts tests/segments/context.test.ts
git commit -m "feat: add context provider and segments (tokens, percent)"
```

---

### Task 11: Git provider and segments

**Files:**
- Create: `src/providers/git.ts`
- Create: `src/segments/git.branch.ts`
- Create: `src/segments/git.insertions.ts`
- Create: `src/segments/git.deletions.ts`
- Create: `tests/providers/git.test.ts`
- Create: `tests/segments/git.test.ts`

- [ ] **Step 1: Write failing provider test**

Create `tests/providers/git.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { gitProvider, gitAvailable } from '../../src/providers/git.js';
import type { SessionData } from '../../src/types.js';

describe('gitAvailable', () => {
  it('returns true for a git repo directory', async () => {
    // Use the ccnow project dir itself (which is a git repo)
    const result = await gitAvailable(process.cwd());
    expect(result).toBe(true);
  });

  it('returns false for /tmp', async () => {
    const result = await gitAvailable('/tmp');
    expect(result).toBe(false);
  });
});

describe('git provider', () => {
  it('resolves branch name for a git repo', async () => {
    const session: SessionData = { cwd: process.cwd() };
    const data = await gitProvider.resolve(session) as any;
    expect(typeof data.branch).toBe('string');
    expect(data.branch.length).toBeGreaterThan(0);
  });

  it('returns null fields for non-git directory', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = await gitProvider.resolve(session) as any;
    expect(data.branch).toBeNull();
    expect(data.insertions).toBeNull();
    expect(data.deletions).toBeNull();
  });

  it('has correct name', () => {
    expect(gitProvider.name).toBe('git');
  });
});
```

- [ ] **Step 2: Write failing segment tests**

Create `tests/segments/git.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { gitBranchSegment } from '../../src/segments/git.branch.js';
import { gitInsertionsSegment } from '../../src/segments/git.insertions.js';
import { gitDeletionsSegment } from '../../src/segments/git.deletions.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('git.branch segment', () => {
  it('returns branch name', () => {
    const provider = { branch: 'main', insertions: 5, deletions: 3 };
    expect(gitBranchSegment.render({ session, provider })).toBe('main');
  });

  it('returns null when branch is null', () => {
    const provider = { branch: null, insertions: null, deletions: null };
    expect(gitBranchSegment.render({ session, provider })).toBeNull();
  });

  it('declares git provider', () => {
    expect(gitBranchSegment.provider).toBe('git');
  });
});

describe('git.insertions segment', () => {
  it('returns insertion count as string', () => {
    const provider = { branch: 'main', insertions: 12, deletions: 3 };
    expect(gitInsertionsSegment.render({ session, provider })).toBe('12');
  });

  it('returns null when insertions is 0', () => {
    const provider = { branch: 'main', insertions: 0, deletions: 0 };
    expect(gitInsertionsSegment.render({ session, provider })).toBeNull();
  });

  it('returns null when insertions is null', () => {
    const provider = { branch: null, insertions: null, deletions: null };
    expect(gitInsertionsSegment.render({ session, provider })).toBeNull();
  });
});

describe('git.deletions segment', () => {
  it('returns deletion count as string', () => {
    const provider = { branch: 'main', insertions: 0, deletions: 7 };
    expect(gitDeletionsSegment.render({ session, provider })).toBe('7');
  });

  it('returns null when deletions is 0', () => {
    const provider = { branch: 'main', insertions: 0, deletions: 0 };
    expect(gitDeletionsSegment.render({ session, provider })).toBeNull();
  });
});
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `npm test -- tests/providers/git.test.ts tests/segments/git.test.ts`
Expected: FAIL

- [ ] **Step 4: Implement git provider**

Create `src/providers/git.ts`:

```ts
import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import type { DataProvider, SessionData } from '../types.js';

const execFileAsync = promisify(execFile);

export interface GitData {
  branch: string | null;
  insertions: number | null;
  deletions: number | null;
}

async function exec(cmd: string, args: string[], cwd: string): Promise<string> {
  try {
    const { stdout } = await execFileAsync(cmd, args, { cwd, timeout: 5000 });
    return stdout.trim();
  } catch {
    return '';
  }
}

export async function gitAvailable(cwd: string): Promise<boolean> {
  const result = await exec('git', ['-C', cwd, 'rev-parse', '--git-dir'], cwd);
  return result !== '';
}

export const gitProvider: DataProvider = {
  name: 'git',
  async resolve(session: SessionData): Promise<GitData> {
    const cwd = session.cwd;
    const isGit = await gitAvailable(cwd);

    if (!isGit) {
      return { branch: null, insertions: null, deletions: null };
    }

    const branch = await exec('git', ['branch', '--show-current'], cwd) || null;

    const diffstat = await exec('git', ['diff', '--shortstat', 'HEAD'], cwd);
    let insertions: number | null = null;
    let deletions: number | null = null;

    if (diffstat) {
      const insMatch = diffstat.match(/(\d+) insertion/);
      const delMatch = diffstat.match(/(\d+) deletion/);
      if (insMatch) insertions = parseInt(insMatch[1], 10);
      if (delMatch) deletions = parseInt(delMatch[1], 10);
    }

    return { branch, insertions, deletions };
  },
};
```

- [ ] **Step 5: Implement git segments**

Create `src/segments/git.branch.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { GitData } from '../providers/git.js';

export const gitBranchSegment: Segment = {
  name: 'git.branch',
  provider: 'git',
  render(context: SegmentContext): string | null {
    const data = context.provider as GitData | undefined;
    return data?.branch ?? null;
  },
};
```

Create `src/segments/git.insertions.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { GitData } from '../providers/git.js';

export const gitInsertionsSegment: Segment = {
  name: 'git.insertions',
  provider: 'git',
  render(context: SegmentContext): string | null {
    const data = context.provider as GitData | undefined;
    if (!data?.insertions) return null;
    return `${data.insertions}`;
  },
};
```

Create `src/segments/git.deletions.ts`:

```ts
import type { Segment, SegmentContext } from '../types.js';
import type { GitData } from '../providers/git.js';

export const gitDeletionsSegment: Segment = {
  name: 'git.deletions',
  provider: 'git',
  render(context: SegmentContext): string | null {
    const data = context.provider as GitData | undefined;
    if (!data?.deletions) return null;
    return `${data.deletions}`;
  },
};
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `npm test -- tests/providers/git.test.ts tests/segments/git.test.ts`
Expected: PASS — all 12 tests pass

- [ ] **Step 7: Commit**

```bash
git add src/providers/git.ts src/segments/git.branch.ts src/segments/git.insertions.ts src/segments/git.deletions.ts tests/providers/git.test.ts tests/segments/git.test.ts
git commit -m "feat: add git provider and segments (branch, insertions, deletions)"
```

---

### Task 12: Segment and provider index files

**Files:**
- Create: `src/segments/index.ts`
- Create: `src/providers/index.ts`

- [ ] **Step 1: Create segment index**

Create `src/segments/index.ts`:

```ts
import type { SegmentRegistry } from '../registry.js';
import { literalSegment } from './literal.js';
import { sepSegment } from './sep.js';
import { pwdNameSegment } from './pwd.name.js';
import { pwdPathSegment } from './pwd.path.js';
import { pwdSmartSegment } from './pwd.smart.js';
import { gitBranchSegment } from './git.branch.js';
import { gitInsertionsSegment } from './git.insertions.js';
import { gitDeletionsSegment } from './git.deletions.js';
import { contextTokensSegment } from './context.tokens.js';
import { contextPercentSegment } from './context.percent.js';

export function registerBuiltinSegments(registry: SegmentRegistry): void {
  registry.register(literalSegment);
  registry.register(sepSegment);
  registry.register(pwdNameSegment);
  registry.register(pwdPathSegment);
  registry.register(pwdSmartSegment);
  registry.register(gitBranchSegment);
  registry.register(gitInsertionsSegment);
  registry.register(gitDeletionsSegment);
  registry.register(contextTokensSegment);
  registry.register(contextPercentSegment);
}
```

- [ ] **Step 2: Create provider index**

Create `src/providers/index.ts`:

```ts
import type { ProviderRegistry } from '../providers.js';
import { pwdProvider } from './pwd.js';
import { gitProvider } from './git.js';
import { contextProvider } from './context.js';

export function registerBuiltinProviders(registry: ProviderRegistry): void {
  registry.register(pwdProvider);
  registry.register(gitProvider);
  registry.register(contextProvider);
}
```

- [ ] **Step 3: Commit**

```bash
git add src/segments/index.ts src/providers/index.ts
git commit -m "feat: add segment and provider registration index files"
```

---

## Chunk 4: DSL, Presets, CLI, and Runner

### Task 13: DSL factory functions

**Files:**
- Create: `src/dsl/index.ts`
- Create: `tests/dsl/index.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/dsl/index.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import {
  StatusLine, Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../../src/dsl/index.js';
import type { SegmentNode } from '../../src/types.js';

describe('DSL factory functions', () => {
  it('Literal creates a literal node', () => {
    const node = Literal({ text: 'hello' });
    expect(node.type).toBe('literal');
    expect(node.props?.text).toBe('hello');
  });

  it('Sep creates a separator node', () => {
    const node = Sep({ char: '|', dim: true });
    expect(node.type).toBe('sep');
    expect(node.props?.char).toBe('|');
    expect(node.style?.dim).toBe(true);
  });

  it('Pwd creates a pwd.smart node by default', () => {
    const node = Pwd({ color: 'cyan' });
    expect(node.type).toBe('pwd.smart');
    expect(node.provider).toBe('pwd');
    expect(node.style?.color).toBe('cyan');
  });

  it('Pwd respects style override', () => {
    const node = Pwd({ style: 'name', color: 'blue' });
    expect(node.type).toBe('pwd.name');
  });

  it('Branch creates a git.branch node', () => {
    const node = Branch({ color: 'white', icon: '\ue0a0 ' });
    expect(node.type).toBe('git.branch');
    expect(node.provider).toBe('git');
    expect(node.style?.icon).toBe('\ue0a0 ');
  });

  it('Git creates a composite node with trailing closure', () => {
    const node = Git()(() => [
      Branch({ color: 'white' }),
      Literal({ text: ' ' }),
    ]);
    expect(node.type).toBe('git');
    expect(node.children).toHaveLength(2);
    expect(node.children![0].type).toBe('git.branch');
  });

  it('Git accepts enabled function', () => {
    const enabledFn = (s: any) => s.cwd !== '';
    const node = Git({ enabled: enabledFn })(() => [Branch()]);
    expect(typeof node.enabled).toBe('function');
  });

  it('Context creates a composite node', () => {
    const node = Context()(() => [
      Tokens({ bold: true }),
      Percent(),
    ]);
    expect(node.type).toBe('context');
    expect(node.children).toHaveLength(2);
  });

  it('StatusLine returns flat array of nodes', () => {
    const tree = StatusLine(() => [
      Pwd({ color: 'cyan' }),
      Sep({ char: '|' }),
    ]);
    expect(tree).toHaveLength(2);
    expect(tree[0].type).toBe('pwd.smart');
    expect(tree[1].type).toBe('sep');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/dsl/`
Expected: FAIL

- [ ] **Step 3: Implement DSL factory functions**

Create `src/dsl/index.ts`:

```ts
import type { SegmentNode, EnabledFn, StyleAttrs } from '../types.js';

interface BaseProps extends Partial<StyleAttrs> {
  enabled?: boolean | EnabledFn;
}

interface PwdProps extends BaseProps {
  style?: 'name' | 'path' | 'smart';
}

interface LiteralProps {
  text: string;
}

interface SepProps extends Partial<StyleAttrs> {
  char?: string;
}

interface CompositeProps extends BaseProps {}

function extractStyle(props: Partial<StyleAttrs>): StyleAttrs | undefined {
  const style: StyleAttrs = {};
  let hasStyle = false;
  for (const key of ['color', 'bold', 'dim', 'italic', 'icon', 'prefix', 'suffix'] as const) {
    if (props[key] !== undefined) {
      (style as any)[key] = props[key];
      hasStyle = true;
    }
  }
  return hasStyle ? style : undefined;
}

export function Literal(props: LiteralProps): SegmentNode {
  return { type: 'literal', props: { text: props.text } };
}

export function Sep(props: SepProps = {}): SegmentNode {
  const { char, ...styleProps } = props;
  return {
    type: 'sep',
    props: char !== undefined ? { char } : undefined,
    style: extractStyle(styleProps),
  };
}

export function Pwd(props: PwdProps = {}): SegmentNode {
  const { style: variant = 'smart', enabled, ...styleProps } = props;
  return {
    type: `pwd.${variant}`,
    provider: 'pwd',
    enabled,
    style: extractStyle(styleProps),
  };
}

export function Branch(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.branch', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function Insertions(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.insertions', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function Deletions(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.deletions', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function Tokens(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.tokens', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function Percent(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.percent', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function Git(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'git',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function Context(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'context',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function StatusLine(children: () => SegmentNode[]): SegmentNode[] {
  return children();
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/dsl/`
Expected: PASS — all 10 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/dsl/index.ts tests/dsl/index.test.ts
git commit -m "feat: add DSL factory functions for composing segment trees"
```

---

### Task 14: Default and minimal presets

**Files:**
- Create: `src/presets/default.ts`
- Create: `src/presets/minimal.ts`
- Create: `src/presets/index.ts`
- Create: `tests/presets/default.test.ts`
- Create: `tests/presets/minimal.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/presets/default.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { defaultPreset } from '../../src/presets/default.js';

describe('default preset', () => {
  it('returns an array of segment nodes', () => {
    expect(Array.isArray(defaultPreset)).toBe(true);
    expect(defaultPreset.length).toBeGreaterThan(0);
  });

  it('starts with pwd.smart', () => {
    expect(defaultPreset[0].type).toBe('pwd.smart');
  });

  it('contains a git composite', () => {
    const git = defaultPreset.find((n) => n.type === 'git');
    expect(git).toBeDefined();
    expect(git?.children?.length).toBeGreaterThan(0);
  });

  it('contains a context composite', () => {
    const ctx = defaultPreset.find((n) => n.type === 'context');
    expect(ctx).toBeDefined();
    expect(ctx?.children?.length).toBeGreaterThan(0);
  });

  it('contains separators', () => {
    const seps = defaultPreset.filter((n) => n.type === 'sep');
    expect(seps.length).toBeGreaterThanOrEqual(2);
  });
});
```

Create `tests/presets/minimal.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { minimalPreset } from '../../src/presets/minimal.js';

describe('minimal preset', () => {
  it('returns an array of segment nodes', () => {
    expect(Array.isArray(minimalPreset)).toBe(true);
  });

  it('contains pwd and git.branch', () => {
    const types = minimalPreset.map((n) => n.type);
    expect(types).toContain('pwd.name');
    expect(types).toContain('git.branch');
  });

  it('is shorter than default preset', () => {
    expect(minimalPreset.length).toBeLessThanOrEqual(4);
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npm test -- tests/presets/`
Expected: FAIL

- [ ] **Step 3: Implement presets**

Create `src/presets/default.ts`:

```ts
import {
  StatusLine, Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const defaultPreset: SegmentNode[] = StatusLine(() => [
  Pwd({ color: 'cyan', bold: true }),
  Sep({ char: '|', dim: true }),
  Git()(() => [
    Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
    Literal({ text: ' [' }),
    Insertions({ color: 'green', prefix: '+' }),
    Literal({ text: ' ' }),
    Deletions({ color: 'red', prefix: '-' }),
    Literal({ text: ']' }),
  ]),
  Sep({ char: '|', dim: true }),
  Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ bold: true }),
    Literal({ text: ' (' }),
    Percent(),
    Literal({ text: ')' }),
  ]),
]);
```

Create `src/presets/minimal.ts`:

```ts
import { StatusLine, Pwd, Sep, Branch } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  Pwd({ style: 'name', color: 'cyan', bold: true }),
  Sep({ char: '|', dim: true }),
  Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
]);
```

Create `src/presets/full.ts`:

```ts
import {
  StatusLine, Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const fullPreset: SegmentNode[] = StatusLine(() => [
  Pwd({ style: 'path', color: 'cyan', bold: true }),
  Sep({ char: '|', dim: true }),
  Git()(() => [
    Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
    Literal({ text: ' [' }),
    Insertions({ color: 'green', prefix: '+' }),
    Literal({ text: ' ' }),
    Deletions({ color: 'red', prefix: '-' }),
    Literal({ text: ']' }),
  ]),
  Sep({ char: '|', dim: true }),
  Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ bold: true }),
    Literal({ text: ' (' }),
    Percent(),
    Literal({ text: ')' }),
  ]),
]);
```

Create `src/presets/index.ts`:

```ts
import type { SegmentNode } from '../types.js';
import { defaultPreset } from './default.js';
import { minimalPreset } from './minimal.js';
import { fullPreset } from './full.js';

const presets = new Map<string, SegmentNode[]>([
  ['default', defaultPreset],
  ['minimal', minimalPreset],
  ['full', fullPreset],
]);

export function getPreset(name: string): SegmentNode[] | undefined {
  return presets.get(name);
}

export function listPresets(): string[] {
  return [...presets.keys()];
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npm test -- tests/presets/`
Expected: PASS — all 8 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/presets/ tests/presets/
git commit -m "feat: add default, minimal, and full presets using DSL"
```

---

### Task 15: CLI argument parser

**Files:**
- Create: `src/cli-parser.ts`
- Create: `tests/cli-parser.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/cli-parser.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { parseArgs, type CliArgs } from '../src/cli-parser.js';

describe('parseArgs', () => {
  it('returns defaults when no args', () => {
    const args = parseArgs([]);
    expect(args.preset).toBe('default');
    expect(args.segments).toHaveLength(0);
    expect(args.config).toBeUndefined();
    expect(args.format).toBe('ansi');
    expect(args.tee).toBeUndefined();
  });

  it('parses --preset flag', () => {
    const args = parseArgs(['--preset', 'minimal']);
    expect(args.preset).toBe('minimal');
  });

  it('parses --preset=value syntax', () => {
    const args = parseArgs(['--preset=minimal']);
    expect(args.preset).toBe('minimal');
  });

  it('parses --config flag', () => {
    const args = parseArgs(['--config', '/path/to/config.json']);
    expect(args.config).toBe('/path/to/config.json');
  });

  it('parses --format flag', () => {
    const args = parseArgs(['--format', 'plain']);
    expect(args.format).toBe('plain');
  });

  it('parses --tee flag', () => {
    const args = parseArgs(['--tee', '/tmp/session.json']);
    expect(args.tee).toBe('/tmp/session.json');
  });

  it('extracts segment flags in order', () => {
    const args = parseArgs(['--pwd', '--sep', '--git', '--sep', '--context']);
    expect(args.segments).toEqual(['pwd', 'sep', 'git', 'sep', 'context']);
  });

  it('preserves duplicate sep flags', () => {
    const args = parseArgs(['--sep', '--sep']);
    expect(args.segments).toEqual(['sep', 'sep']);
  });

  it('deduplicates non-sep composite flags', () => {
    const args = parseArgs(['--git', '--git']);
    expect(args.segments).toEqual(['git']);
  });

  it('mixes segment and value flags', () => {
    const args = parseArgs(['--pwd', '--preset', 'minimal', '--sep', '--git']);
    expect(args.segments).toEqual(['pwd', 'sep', 'git']);
    expect(args.preset).toBe('minimal');
  });

  it('parses --help flag', () => {
    const args = parseArgs(['--help']);
    expect(args.help).toBe(true);
  });

  it('parses --version flag', () => {
    const args = parseArgs(['--version']);
    expect(args.version).toBe(true);
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/cli-parser.test.ts`
Expected: FAIL

- [ ] **Step 3: Implement CLI parser**

Create `src/cli-parser.ts`:

```ts
export interface CliArgs {
  preset: string;
  segments: string[];
  config?: string;
  format: 'ansi' | 'plain';
  tee?: string;
  help?: boolean;
  version?: boolean;
}

const SEGMENT_FLAGS = new Set(['pwd', 'git', 'context', 'sep']);
const VALUE_FLAGS = new Set(['preset', 'config', 'format', 'tee']);

export function parseArgs(argv: string[]): CliArgs {
  const result: CliArgs = {
    preset: 'default',
    segments: [],
    format: 'ansi',
  };

  const seenComposites = new Set<string>();
  let i = 0;

  while (i < argv.length) {
    const arg = argv[i];

    if (arg === '--help') {
      result.help = true;
      i++;
      continue;
    }

    if (arg === '--version') {
      result.version = true;
      i++;
      continue;
    }

    // Handle --key=value syntax
    const eqMatch = arg.match(/^--(\w[\w-]*)=(.+)$/);
    if (eqMatch) {
      const [, key, value] = eqMatch;
      if (VALUE_FLAGS.has(key)) {
        (result as any)[key] = value;
      }
      i++;
      continue;
    }

    // Handle --key value syntax for value flags
    const flagMatch = arg.match(/^--(\w[\w-]*)$/);
    if (flagMatch) {
      const key = flagMatch[1];

      if (VALUE_FLAGS.has(key)) {
        const value = argv[i + 1];
        if (value !== undefined) {
          (result as any)[key] = value;
          i += 2;
          continue;
        }
      }

      if (SEGMENT_FLAGS.has(key)) {
        // sep is always allowed as duplicate
        if (key === 'sep' || !seenComposites.has(key)) {
          result.segments.push(key);
          if (key !== 'sep') seenComposites.add(key);
        }
      }

      i++;
      continue;
    }

    i++;
  }

  return result;
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/cli-parser.test.ts`
Expected: PASS — all 12 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/cli-parser.ts tests/cli-parser.test.ts
git commit -m "feat: add CLI argument parser with ordered segment flags"
```

---

### Task 16: JSON config loader

**Files:**
- Create: `src/config.ts`
- Create: `tests/config.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/config.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { parseConfig } from '../src/config.js';
import type { SegmentNode } from '../src/types.js';

describe('parseConfig', () => {
  it('parses a simple segment list', () => {
    const json = {
      segments: [
        { segment: 'pwd.smart', style: { color: 'cyan' } },
        { segment: 'sep', props: { char: '|' } },
      ],
    };
    const tree = parseConfig(json);
    expect(tree).toHaveLength(2);
    expect(tree[0].type).toBe('pwd.smart');
    expect(tree[0].style?.color).toBe('cyan');
    expect(tree[1].type).toBe('sep');
    expect(tree[1].props?.char).toBe('|');
  });

  it('parses composite segments with children', () => {
    const json = {
      segments: [
        {
          segment: 'git',
          children: [
            { segment: 'git.branch', style: { color: 'white' } },
          ],
        },
      ],
    };
    const tree = parseConfig(json);
    expect(tree[0].type).toBe('git');
    expect(tree[0].children).toHaveLength(1);
    expect(tree[0].children![0].type).toBe('git.branch');
  });

  it('handles enabled boolean', () => {
    const json = {
      segments: [
        { segment: 'git', enabled: false, children: [] },
      ],
    };
    const tree = parseConfig(json);
    expect(tree[0].enabled).toBe(false);
  });

  it('returns empty array for missing segments', () => {
    expect(parseConfig({})).toEqual([]);
    expect(parseConfig({ segments: [] })).toEqual([]);
  });

  it('maps segment field to type and infers provider', () => {
    const json = {
      segments: [{ segment: 'git.branch' }],
    };
    const tree = parseConfig(json);
    expect(tree[0].type).toBe('git.branch');
    expect(tree[0].provider).toBe('git');
  });

  it('does not infer provider for literal and sep', () => {
    const json = {
      segments: [
        { segment: 'literal', props: { text: 'hi' } },
        { segment: 'sep' },
      ],
    };
    const tree = parseConfig(json);
    expect(tree[0].provider).toBeUndefined();
    expect(tree[1].provider).toBeUndefined();
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/config.test.ts`
Expected: FAIL

- [ ] **Step 3: Implement config parser**

Create `src/config.ts`:

```ts
import type { SegmentNode } from './types.js';

interface JsonSegment {
  segment: string;
  style?: Record<string, unknown>;
  props?: Record<string, unknown>;
  enabled?: boolean;
  children?: JsonSegment[];
}

interface JsonConfig {
  segments?: JsonSegment[];
}

const NO_PROVIDER = new Set(['literal', 'sep']);

function inferProvider(type: string): string | undefined {
  if (NO_PROVIDER.has(type)) return undefined;
  const dotIndex = type.indexOf('.');
  if (dotIndex > 0) return type.slice(0, dotIndex);
  return undefined;
}

function mapSegment(json: JsonSegment): SegmentNode {
  const node: SegmentNode = {
    type: json.segment,
  };

  const provider = inferProvider(json.segment);
  if (provider) node.provider = provider;

  if (json.style) node.style = json.style as any;
  if (json.props) node.props = json.props;
  if (json.enabled !== undefined) node.enabled = json.enabled;
  if (json.children) {
    node.children = json.children.map(mapSegment);
  }

  return node;
}

export function parseConfig(config: Record<string, unknown>): SegmentNode[] {
  const typed = config as JsonConfig;
  if (!typed.segments || !Array.isArray(typed.segments)) return [];
  return typed.segments.map(mapSegment);
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm test -- tests/config.test.ts`
Expected: PASS — all 6 tests pass

- [ ] **Step 5: Commit**

```bash
git add src/config.ts tests/config.test.ts
git commit -m "feat: add JSON config parser with provider inference"
```

---

### Task 17: Runner pipeline

**Files:**
- Create: `src/runner.ts`
- Create: `tests/runner.test.ts`

- [ ] **Step 1: Write failing tests**

Create `tests/runner.test.ts`:

```ts
import { describe, it, expect } from '@jest/globals';
import { run } from '../src/runner.js';

describe('run', () => {
  it('renders default preset with basic session data', async () => {
    const input = JSON.stringify({
      cwd: process.cwd(),
      context_window: {
        used_percentage: 42,
        current_usage: { input_tokens: 38000, cache_creation_input_tokens: 2000, cache_read_input_tokens: 1500 },
      },
    });
    const output = await run({ preset: 'default', segments: [], format: 'ansi' }, input);
    expect(output.length).toBeGreaterThan(0);
    // Should contain context info
    expect(output).toContain('42K');
    expect(output).toContain('42%');
  });

  it('renders plain format without ANSI codes', async () => {
    const input = JSON.stringify({ cwd: '/tmp' });
    const output = await run({ preset: 'minimal', segments: [], format: 'plain' }, input);
    // Plain format should not contain ANSI escape codes
    expect(output).not.toMatch(/\x1b\[/);
  });

  it('returns empty string for invalid stdin', async () => {
    const output = await run({ preset: 'default', segments: [], format: 'ansi' }, 'not json');
    expect(output).toBe('');
  });

  it('uses segment flags when provided', async () => {
    const input = JSON.stringify({
      cwd: '/tmp',
      context_window: { used_percentage: 10, current_usage: { input_tokens: 5000, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 } },
    });
    const output = await run({ preset: 'default', segments: ['context'], format: 'plain' }, input);
    expect(output).toContain('5K');
    // Should not contain pwd since we only requested context
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm test -- tests/runner.test.ts`
Expected: FAIL

- [ ] **Step 3: Implement runner**

Create `src/runner.ts`:

```ts
import chalk from 'chalk';
import { readFileSync } from 'node:fs';
import type { CliArgs } from './cli-parser.js';
import type { SegmentNode } from './types.js';
import { parseSession } from './session.js';
import { parseConfig } from './config.js';
import { SegmentRegistry } from './registry.js';
import { ProviderRegistry } from './providers.js';
import { registerBuiltinSegments } from './segments/index.js';
import { registerBuiltinProviders } from './providers/index.js';
import { renderTree } from './render.js';
import { getPreset } from './presets/index.js';
import { buildCompositeTree } from './composites.js';

export async function run(args: CliArgs, stdin: string): Promise<string> {
  // Parse session
  const session = parseSession(stdin);
  if (!session) return '';

  // Set chalk level based on format (save and restore to avoid global mutation)
  const originalLevel = chalk.level;
  if (args.format === 'plain') {
    chalk.level = 0;
  }

  try {
    // Build registries
    const segmentRegistry = new SegmentRegistry();
    registerBuiltinSegments(segmentRegistry);

    const providerRegistry = new ProviderRegistry();
    registerBuiltinProviders(providerRegistry);

    // Resolve render tree: config file > CLI segment flags > preset
    let tree: SegmentNode[];

    if (args.config) {
      try {
        const configJson = JSON.parse(readFileSync(args.config, 'utf-8'));
        const configTree = parseConfig(configJson);
        tree = configTree.length > 0 ? configTree : getPreset(args.preset) ?? getPreset('default')!;
      } catch (err) {
        process.stderr.write(`ccnow: failed to load config: ${err}\n`);
        tree = getPreset(args.preset) ?? getPreset('default')!;
      }
    } else if (args.segments.length > 0) {
      tree = buildCompositeTree(args.segments);
    } else {
      tree = getPreset(args.preset) ?? getPreset('default')!;
    }

    // Resolve providers
    const providerNames = providerRegistry.collectProviderNames(tree);
    const providerData = await providerRegistry.resolveAll([...providerNames], session);

    // Render
    return renderTree(tree, segmentRegistry, session, providerData);
  } finally {
    chalk.level = originalLevel;
  }
}
```

- [ ] **Step 4: Create composites helper**

This maps CLI segment flags to their composite render trees.

Create `src/composites.ts`:

```ts
import {
  Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from './dsl/index.js';
import type { SegmentNode } from './types.js';

const compositeBuilders: Record<string, () => SegmentNode> = {
  pwd: () => Pwd({ color: 'cyan', bold: true }),
  sep: () => Sep({ char: '|', dim: true }),
  git: () => Git()(() => [
    Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
    Literal({ text: ' [' }),
    Insertions({ color: 'green', prefix: '+' }),
    Literal({ text: ' ' }),
    Deletions({ color: 'red', prefix: '-' }),
    Literal({ text: ']' }),
  ]),
  context: () => Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ bold: true }),
    Literal({ text: ' (' }),
    Percent(),
    Literal({ text: ')' }),
  ]),
};

export function buildCompositeTree(segments: string[]): SegmentNode[] {
  return segments
    .map((name) => compositeBuilders[name]?.())
    .filter((node): node is SegmentNode => node !== undefined);
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `npm test -- tests/runner.test.ts`
Expected: PASS — all 4 tests pass

- [ ] **Step 6: Commit**

```bash
git add src/runner.ts src/composites.ts tests/runner.test.ts
git commit -m "feat: add runner pipeline with composite tree builder"
```

---

### Task 18: CLI entry point

**Files:**
- Create: `src/cli.ts`

- [ ] **Step 1: Implement CLI entry point**

Create `src/cli.ts`:

```ts
#!/usr/bin/env node

import { readFileSync, writeFileSync } from 'node:fs';
import { parseArgs } from './cli-parser.js';
import { run } from './runner.js';

const HELP = `Usage: ccnow [options]

Composable statusline for Claude Code.
Reads session JSON from stdin, outputs styled statusline to stdout.

Options:
  --preset <name>     Use a named preset (default, minimal, full)
  --config <path>     Load JSON config file
  --pwd               Enable pwd composite segment
  --git               Enable git composite segment
  --context           Enable context composite segment
  --sep               Insert separator segment
  --format <type>     Output format: ansi (default), plain
  --tee <path>        Write raw stdin JSON to file before processing
  --help              Show help
  --version           Show version

Examples:
  npx -y ccnow
  npx -y ccnow --preset=minimal
  npx -y ccnow --pwd --sep --git --sep --context
  npx -y ccnow --config ~/.claude/ccnow.json
`;

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));

  if (args.help) {
    process.stdout.write(HELP);
    return;
  }

  if (args.version) {
    // Read version from package.json at runtime
    try {
      const pkg = JSON.parse(readFileSync(new URL('../package.json', import.meta.url), 'utf-8'));
      process.stdout.write(`${pkg.version}\n`);
    } catch {
      process.stdout.write('unknown\n');
    }
    return;
  }

  // Read stdin
  let stdin: string;
  try {
    stdin = readFileSync(0, 'utf-8');
  } catch {
    stdin = '';
  }

  // Tee: write raw stdin to file before processing
  if (args.tee) {
    try {
      writeFileSync(args.tee, stdin, 'utf-8');
    } catch (err) {
      process.stderr.write(`ccnow: failed to write tee file: ${err}\n`);
    }
  }

  const output = await run(args, stdin);
  if (output) process.stdout.write(output);
}

main().catch((err) => {
  process.stderr.write(`ccnow: ${err}\n`);
  process.exit(1);
});
```

- [ ] **Step 2: Build and test manually**

Run:
```bash
cd /Users/jheddings/Projects/ccnow && npx tsc
echo '{"cwd":"/Users/jheddings/Projects/ccnow","context_window":{"used_percentage":42,"current_usage":{"input_tokens":38000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":1500}}}' | node dist/cli.js
```
Expected: styled statusline output with pwd, git info, and context stats

- [ ] **Step 3: Test --help flag**

Run: `node dist/cli.js --help`
Expected: help text output

- [ ] **Step 4: Test --preset=minimal**

Run:
```bash
echo '{"cwd":"/Users/jheddings/Projects/ccnow"}' | node dist/cli.js --preset=minimal
```
Expected: shorter output with just pwd and branch

- [ ] **Step 5: Test --format=plain**

Run:
```bash
echo '{"cwd":"/Users/jheddings/Projects/ccnow","context_window":{"used_percentage":42,"current_usage":{"input_tokens":38000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":1500}}}' | node dist/cli.js --format=plain
```
Expected: same content but no ANSI escape codes

- [ ] **Step 6: Commit**

```bash
git add src/cli.ts
git commit -m "feat: add CLI entry point with stdin/stdout pipeline"
```

---

### Task 19: Run full test suite and verify

- [ ] **Step 1: Run all tests**

Run: `cd /Users/jheddings/Projects/ccnow && npm test`
Expected: ALL tests pass

- [ ] **Step 2: Fix any failures**

If any tests fail, fix the issues and re-run.

- [ ] **Step 3: Final manual integration test**

Run:
```bash
echo '{"cwd":"/Users/jheddings/Projects/ccnow","context_window":{"used_percentage":42,"current_usage":{"input_tokens":38000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":1500}}}' | node dist/cli.js
```
Expected: output matches the style of the existing `~/.claude/statusline.sh`

- [ ] **Step 4: Commit any fixes**

```bash
git add -A && git commit -m "fix: test suite fixes"
```

---

### Task 20: Wire up Claude Code statusline

- [ ] **Step 1: Build the project**

Run: `cd /Users/jheddings/Projects/ccnow && npx tsc`

- [ ] **Step 2: Test with absolute path in settings**

Temporarily update `~/.claude/settings.json` to use ccnow:

```json
{
  "statusLine": {
    "type": "command",
    "command": "node /Users/jheddings/Projects/ccnow/dist/cli.js",
    "padding": 0
  }
}
```

- [ ] **Step 3: Verify in Claude Code**

Open a new Claude Code session and confirm the statusline renders correctly.

- [ ] **Step 4: Restore original settings**

Restore the original `statusline.sh` command in settings if needed, or keep ccnow if it works well.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "chore: final integration verification"
```
