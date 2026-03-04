import { request } from '../api/client.js';

export function publishArticle(title, content) {
  return request('/articles', {
    method: 'POST',
    data: { title, content },
  });
}

export function recommendFeed(limit = 20) {
  return request(`/articles/recommend?limit=${encodeURIComponent(limit)}`);
}

export function followingFeed() {
  return request('/articles/feed');
}

export function articleDetail(id) {
  return request(`/articles/${id}`);
}

export function likeArticle(id) {
  return request(`/interactions/articles/${id}/like`, { method: 'POST' });
}

export function unlikeArticle(id) {
  return request(`/interactions/articles/${id}/like`, { method: 'DELETE' });
}

export function addComment(id, content) {
  return request(`/interactions/articles/${id}/comments`, {
    method: 'POST',
    data: { content },
  });
}
