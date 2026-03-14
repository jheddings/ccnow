import type { StyleAttrs } from './types.js';

const RESET = '\x1b[0m';

const COLORS: Record<string, string> = {
  black: '\x1b[30m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
  white: '\x1b[37m',
  blackBright: '\x1b[90m',
  redBright: '\x1b[91m',
  greenBright: '\x1b[92m',
  yellowBright: '\x1b[93m',
  blueBright: '\x1b[94m',
  magentaBright: '\x1b[95m',
  cyanBright: '\x1b[96m',
  whiteBright: '\x1b[97m',
};

const BOLD = '\x1b[1m';
const DIM = '\x1b[2m';
const ITALIC = '\x1b[3m';

let colorsEnabled = true;

export function setColorLevel(level: number): void {
  colorsEnabled = level > 0;
}

function resolveColor(color: string): string {
  // Named color: lookup in table
  const named = COLORS[color];
  if (named) return named;

  // Hex color: #RRGGBB → truecolor \e[38;2;R;G;Bm
  if (color.startsWith('#') && color.length === 7) {
    const r = parseInt(color.slice(1, 3), 16);
    const g = parseInt(color.slice(3, 5), 16);
    const b = parseInt(color.slice(5, 7), 16);
    if (!isNaN(r) && !isNaN(g) && !isNaN(b)) {
      return `\x1b[38;2;${r};${g};${b}m`;
    }
  }

  // Numeric string: 256-color index → \e[38;5;Nm
  const num = parseInt(color, 10);
  if (!isNaN(num) && num >= 0 && num <= 255) {
    return `\x1b[38;5;${num}m`;
  }

  return '';
}

export function applyStyle(value: string, style: StyleAttrs | undefined): string {
  if (!style) return value;

  // Build the decorated string: prefix + value + suffix
  let result = value;
  if (style.prefix) result = style.prefix + result;
  if (style.suffix) result = result + style.suffix;

  if (colorsEnabled) {
    // Build ANSI open sequence: reset first, then modifiers, then color
    let open = '';
    if (style.bold) open += BOLD;
    if (style.dim) open += DIM;
    if (style.italic) open += ITALIC;

    if (style.color) {
      open += resolveColor(style.color);
    }

    if (open) {
      result = RESET + open + result + RESET;
    }
  }

  // Prepend icon outside styled region so it renders in default color
  if (style.icon) result = style.icon + result;

  return result;
}
