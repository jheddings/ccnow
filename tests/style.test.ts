import { describe, it, expect } from '@jest/globals';
import { applyStyle, setColorLevel } from '../src/style.js';
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
  });

  it('applies suffix after value', () => {
    const result = applyStyle('42', { suffix: '%' });
    expect(result).toContain('42');
    expect(result).toContain('%');
  });

  it('applies icon outside styled region', () => {
    const result = applyStyle('main', { icon: '\ue0a0 ', color: 'white' });
    // Icon should come before the ANSI codes
    expect(result.indexOf('\ue0a0')).toBeLessThan(result.indexOf('\x1b['));
  });

  it('applies named color as ANSI code', () => {
    const result = applyStyle('hello', { color: 'cyan' });
    expect(result).toContain('\x1b[36m');
    expect(result).toContain('hello');
    expect(result).toContain('\x1b[0m');
  });

  it('applies bold as ANSI code', () => {
    const result = applyStyle('hello', { bold: true });
    expect(result).toContain('\x1b[1m');
    expect(result).toContain('hello');
  });

  it('applies italic as ANSI code', () => {
    const result = applyStyle('hello', { italic: true });
    expect(result).toContain('\x1b[3m');
  });

  it('wraps styled segments with reset on both sides', () => {
    const result = applyStyle('hello', { color: 'red' });
    // Should start with reset, have color, have text, end with reset
    expect(result).toBe('\x1b[0m\x1b[31mhello\x1b[0m');
  });


  it('combines multiple style attrs', () => {
    const style: StyleAttrs = { color: 'green', bold: true, prefix: '+' };
    const result = applyStyle('12', style);
    expect(result).toContain('\x1b[1m');  // bold
    expect(result).toContain('\x1b[32m'); // green
    expect(result).toContain('+12');
  });

  it('handles null/undefined style fields gracefully', () => {
    const style: StyleAttrs = { color: undefined, bold: undefined };
    expect(applyStyle('hello', style)).toBe('hello');
  });

  it('resolves 256-color numeric string', () => {
    const result = applyStyle('hello', { color: '39' });
    expect(result).toContain('\x1b[38;5;39m');
  });

  it('resolves hex color to truecolor', () => {
    const result = applyStyle('hello', { color: '#00afff' });
    expect(result).toContain('\x1b[38;2;0;175;255m');
  });

  it('ignores invalid color values', () => {
    const result = applyStyle('hello', { color: 'notacolor' });
    expect(result).toBe('hello');
  });

  it('disables colors when level is 0', () => {
    setColorLevel(0);
    const result = applyStyle('hello', { color: 'red', bold: true });
    expect(result).toBe('hello');
    expect(result).not.toContain('\x1b[');
    setColorLevel(1); // restore
  });
});
