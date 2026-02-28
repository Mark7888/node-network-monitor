import { IApiClient } from './IApiClient';
import { realApiClient } from './realApiClient';

/**
 * The active API client singleton.
 *
 * In normal mode this is the real HTTP implementation.
 * In mock mode (VITE_MOCK_MODE=true) initialiseApiClient() replaces this
 * with the mock implementation that operates on in-memory state.
 *
 * All service files import from here, so the switch is transparent to the rest
 * of the application.
 */
let _client: IApiClient = realApiClient;

export function getApiClient(): IApiClient {
  return _client;
}

/** Called once from main.tsx before the React tree mounts. */
export async function initialiseApiClient(): Promise<void> {
  if (import.meta.env.VITE_MOCK_MODE === 'true') {
    const { mockApiClient, initializeMockClient } = await import('./mockApiClient');
    await initializeMockClient();
    _client = mockApiClient;
  }
}

/** Convenience proxy — individual service files can do: apiClient.getNodes() */
export const apiClient: IApiClient = new Proxy({} as IApiClient, {
  get(_target, prop) {
    return (...args: unknown[]) => {
      const method = _client[prop as keyof IApiClient];
      if (typeof method !== 'function') {
        throw new Error(`ApiClient: "${String(prop)}" is not a valid method`);
      }
      return (method as (...a: unknown[]) => unknown)(...args);
    };
  },
});
