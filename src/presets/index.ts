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
