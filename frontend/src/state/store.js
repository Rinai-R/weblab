const TOKEN_KEY = 'weblab_token';
const USER_KEY = 'weblab_user';
const API_BASE_KEY = 'weblab_api_base';

const defaultApiBase = 'http://localhost:8080/api/v1';

const savedUser = localStorage.getItem(USER_KEY);

export const state = {
  token: localStorage.getItem(TOKEN_KEY) || '',
  user: savedUser ? JSON.parse(savedUser) : null,
  apiBase: localStorage.getItem(API_BASE_KEY) || defaultApiBase,
  activeFeedTab: 'recommend',
  keyword: '',
  recommendFeed: [],
  followingFeed: [],
  articleDetail: null,
  discoverUsers: [],
  mutualUsers: [],
  activePeer: null,
  messages: [],
  loading: false,
};

export function setToken(token) {
  state.token = token || '';
  if (state.token) {
    localStorage.setItem(TOKEN_KEY, state.token);
    return;
  }
  localStorage.removeItem(TOKEN_KEY);
}

export function setUser(user) {
  state.user = user || null;
  if (state.user) {
    localStorage.setItem(USER_KEY, JSON.stringify(state.user));
    return;
  }
  localStorage.removeItem(USER_KEY);
}

export function clearAuthState() {
  setToken('');
  setUser(null);
  state.articleDetail = null;
  state.messages = [];
}

export function setApiBase(base) {
  state.apiBase = (base || '').trim() || defaultApiBase;
  localStorage.setItem(API_BASE_KEY, state.apiBase);
}
