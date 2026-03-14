import { Chalk } from 'chalk';
import type { StyleAttrs } from './types.js';

// Use a chalk instance with forced color level so ANSI codes are always emitted
let chalk = new Chalk({ level: 3 });

export function setColorLevel(level: 0 | 1 | 2 | 3): void {
  chalk = new Chalk({ level });
}

export function applyStyle(value: string, style: StyleAttrs | undefined): string {
  if (!style) return value;

  // Build the decorated string: icon + prefix + value + suffix
  let result = value;
  if (style.prefix) result = style.prefix + result;
  if (style.suffix) result = result + style.suffix;
  if (style.icon) result = style.icon + result;

  // Apply chalk styling to the full string
  let painter: typeof chalk = chalk;

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
