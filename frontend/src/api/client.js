import { clearAuthState, state } from '../state/store.js';

async function readResponse(res) {
  const contentType = res.headers.get('content-type') || '';
  if (contentType.includes('application/json')) {
    return res.json();
  }
  const text = await res.text();
  return { code: res.status, message: text || 'request failed' };
}

export async function request(path, { method = 'GET', data, auth = true } = {}) {
  const headers = { 'Content-Type': 'application/json' };
  if (auth && state.token) {
    headers.Authorization = `Bearer ${state.token}`;
  }

  const res = await fetch(`${state.apiBase}${path}`, {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
  });

  const payload = await readResponse(res);
  if (!res.ok || payload.code !== 0) {
    const message = payload.message || `HTTP ${res.status}`;
    if (res.status === 401) {
      clearAuthState();
    }
    throw new Error(message);
  }

  return payload.data;
}
