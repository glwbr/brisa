export const customFetch = async <T>(
  config: {
    url: string;
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
    headers?: Record<string, string>;
    data?: unknown;
    signal?: AbortSignal;
  },
  options?: RequestInit,
): Promise<T> => {
  const baseUrl = 'http://localhost:8080/api';
  const fullUrl = config.url.startsWith('http') ? config.url : `${baseUrl}${config.url}`;

  const headers = new Headers(config.headers);
  if (!headers.has('Content-Type') && config.data && !(config.data instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(fullUrl, {
    ...options,
    method: config.method,
    headers,
    signal: config.signal,
    body: config.data ? JSON.stringify(config.data) : undefined,
  });

  if (!response.ok) {
    const errorBody = await response.text();
    throw new Error(errorBody || `HTTP error! status: ${response.status}`);
  }

  if (response.status === 204) {
    return {} as T;
  }

  try {
    return await response.json();
  } catch (error) {
    return {} as T;
  }
};

export default customFetch;
