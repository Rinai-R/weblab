import { request } from '../api/client.js';

export function register(username, password) {
  return request('/auth/register', {
    method: 'POST',
    auth: false,
    data: { username, password },
  });
}

export function login(username, password) {
  return request('/auth/login', {
    method: 'POST',
    auth: false,
    data: { username, password },
  });
}

export function me() {
  return request('/auth/me');
}
