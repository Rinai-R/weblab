import { request } from '../api/client.js';

export function discoverUsers() {
  return request('/social/discover');
}

export function followUser(userID) {
  return request(`/social/follow/${userID}`, { method: 'POST' });
}

export function mutualUsers() {
  return request('/social/mutuals');
}

export function conversation(userID) {
  return request(`/social/messages/${userID}`);
}

export function sendMessage(toUserID, content) {
  return request('/social/messages', {
    method: 'POST',
    data: { to_user_id: toUserID, content },
  });
}
