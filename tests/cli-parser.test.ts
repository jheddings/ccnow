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
