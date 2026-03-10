const TOKEN_KEY = 'weblab_token';
const USER_KEY = 'weblab_user';
const API_BASE_KEY = 'weblab_api_base';

const defaultApiBase = 'http://localhost:8080/api/v1';

const savedUser = localStorage.getItem(USER_KEY);

export const state = {
  token: localStorage.getItem(TOKEN_KEY) || '',
  user: savedUser ? JSON.parse(savedUser) : null,
  apiBase: localStorage.getItem(API_BASE_KEY) || defaultApiBase,
  activeStream: 'articles_recommend',
  keyword: '',
  streams: {
    articles_recommend: { items: [], cursor: 0, hasMore: true },
    articles_following: { items: [], cursor: 0, hasMore: true },
    questions_recommend: { items: [], cursor: 0, hasMore: true },
    questions_following: { items: [], cursor: 0, hasMore: true },
  },
  articleDetail: null,
  questionDetail: null,
  answerCursor: 0,
  answerHasMore: true,
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
  state.activeStream = 'articles_recommend';
  state.streams = {
    articles_recommend: { items: [], cursor: 0, hasMore: true },
    articles_following: { items: [], cursor: 0, hasMore: true },
    questions_recommend: { items: [], cursor: 0, hasMore: true },
    questions_following: { items: [], cursor: 0, hasMore: true },
  };
  state.articleDetail = null;
  state.questionDetail = null;
  state.answerCursor = 0;
  state.answerHasMore = true;
  state.discoverUsers = [];
  state.mutualUsers = [];
  state.activePeer = null;
  state.messages = [];
}

export function setApiBase(base) {
  state.apiBase = (base || '').trim() || defaultApiBase;
  localStorage.setItem(API_BASE_KEY, state.apiBase);
}
