import { IApiClient } from './IApiClient';

/**
 * The active API client singleton.
 *
 * In normal mode this is the real HTTP implementation.
 * In mock mode (VITE_MOCK_MODE=true) initialiseApiClient() sets this
 * to the mock implementation that operates on in-memory state.
 *
 * Both branches use dynamic imports so Rollup can tree-shake the
 * real API client (axios, axiosConfig, etc.) out of mock/demo bundles
 * and the mock client out of production bundles.
 *
 * All service files import from here, so the switch is transparent to
 * the rest of the application.
 */
let _client: IApiClient | undefined;

export function getApiClient(): IApiClient {
  if (!_client) throw new Error('ApiClient has not been initialised yet');
  return _client;
}

/** Called once from main.tsx before the React tree mounts. */
export async function initialiseApiClient(): Promise<void> {
  if (import.meta.env.VITE_MOCK_MODE === 'true') {
    const { mockApiClient, initializeMockClient } = await import('./mockApiClient');
    await initializeMockClient();
    _client = mockApiClient;
  } else {
    // Dynamic import keeps axios/realApiClient out of the mock bundle.
    const { realApiClient } = await import('./realApiClient');
    _client = realApiClient;
  }
}

/** Convenience proxy — individual service files can do: apiClient.getNodes() */
export const apiClient: IApiClient = new Proxy({} as IApiClient, {
  get(_target, prop) {
    return (...args: unknown[]) => {
      if (!_client) throw new Error('ApiClient has not been initialised yet');
      const method = _client[prop as keyof IApiClient];
      if (typeof method !== 'function') {
        throw new Error(`ApiClient: "${String(prop)}" is not a valid method`);
      }
      return (method as (...a: unknown[]) => unknown)(...args);
    };
  },
});
