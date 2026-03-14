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
