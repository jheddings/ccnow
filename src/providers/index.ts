import type { ProviderRegistry } from '../providers.js';
import { pwdProvider } from './pwd.js';
import { gitProvider } from './git.js';
import { contextProvider } from './context.js';

export function registerBuiltinProviders(registry: ProviderRegistry): void {
  registry.register(pwdProvider);
  registry.register(gitProvider);
  registry.register(contextProvider);
}
