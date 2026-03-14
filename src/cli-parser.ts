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
