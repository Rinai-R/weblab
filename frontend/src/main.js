import { clearAuthState, setApiBase, setToken, setUser, state } from './state/store.js';
import { login, me, register } from './modules/auth.js';
import {
  addComment,
  articleDetail,
  followingFeed,
  likeArticle,
  publishArticle,
  recommendFeed,
  unlikeArticle,
} from './modules/articles.js';
import {
  conversation,
  discoverUsers,
  followUser,
  mutualUsers,
  sendMessage,
} from './modules/social.js';

const els = {
  notice: document.getElementById('notice'),
  apiBaseInput: document.getElementById('api-base-input'),
  apiBaseSave: document.getElementById('api-base-save'),
  authPanel: document.getElementById('auth-panel'),
  publishForm: document.getElementById('publish-form'),
  feedTabs: document.getElementById('feed-tabs'),
  searchInput: document.getElementById('feed-search-input'),
  refreshBtn: document.getElementById('feed-refresh-btn'),
  feedList: document.getElementById('feed-list'),
  detailPanel: document.getElementById('detail-panel'),
  discoverList: document.getElementById('discover-list'),
  mutualList: document.getElementById('mutual-list'),
  chatPanel: document.getElementById('chat-panel'),
  chatForm: document.getElementById('chat-form'),
  logoutBtn: document.getElementById('logout-btn'),
};

let authMode = 'login';
const CHAT_POLL_INTERVAL_MS = 1000;
let chatPollTimer = null;
let chatPollInFlight = false;

function escapeHTML(input) {
  const text = String(input ?? '');
  return text
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function timeText(value) {
  try {
    return new Date(value).toLocaleString('zh-CN', { hour12: false });
  } catch {
    return String(value ?? '');
  }
}

function notify(message, type = 'info') {
  els.notice.textContent = message;
  els.notice.className = `notice notice-${type}`;
  window.clearTimeout(notify.timer);
  notify.timer = window.setTimeout(() => {
    els.notice.className = 'notice';
    els.notice.textContent = '';
  }, 2600);
}

function stopConversationPolling() {
  if (chatPollTimer) {
    window.clearTimeout(chatPollTimer);
    chatPollTimer = null;
  }
  chatPollInFlight = false;
}

function hasConversationChanged(nextMessages) {
  if (state.messages.length !== nextMessages.length) {
    return true;
  }
  if (nextMessages.length === 0) {
    return false;
  }

  const prevLast = state.messages[state.messages.length - 1];
  const nextLast = nextMessages[nextMessages.length - 1];
  return prevLast.id !== nextLast.id;
}

function scheduleConversationPolling(delay = CHAT_POLL_INTERVAL_MS) {
  if (!state.token || !state.user || !state.activePeer || document.hidden) {
    return;
  }

  if (chatPollTimer) {
    window.clearTimeout(chatPollTimer);
  }

  chatPollTimer = window.setTimeout(async () => {
    if (!state.token || !state.user || !state.activePeer || document.hidden) {
      stopConversationPolling();
      return;
    }
    if (chatPollInFlight) {
      scheduleConversationPolling();
      return;
    }

    const peerID = state.activePeer.id;
    chatPollInFlight = true;
    try {
      const nextMessages = await conversation(peerID);
      if (state.activePeer && state.activePeer.id === peerID && hasConversationChanged(nextMessages)) {
        state.messages = nextMessages;
        renderChatPanel();
      }
    } catch {
      // Keep polling for transient failures; auth failures are handled in request().
      if (!state.token || !state.user) {
        stopConversationPolling();
        renderAll();
        return;
      }
    } finally {
      chatPollInFlight = false;
      if (state.activePeer && state.activePeer.id === peerID) {
        scheduleConversationPolling();
      }
    }
  }, delay);
}

function startConversationPolling() {
  stopConversationPolling();
  scheduleConversationPolling();
}

function requireLogin() {
  if (state.token && state.user) {
    return true;
  }
  notify('请先登录后再操作', 'warn');
  return false;
}

function activeFeed() {
  return state.activeFeedTab === 'following' ? state.followingFeed : state.recommendFeed;
}

function renderAuthPanel() {
  if (state.user) {
    els.authPanel.innerHTML = `
      <div class="user-head">
        <div class="avatar">${escapeHTML(state.user.username.slice(0, 1).toUpperCase())}</div>
        <div>
          <div class="label">当前账号</div>
          <div class="value">${escapeHTML(state.user.username)}</div>
          <div class="muted">UID ${state.user.id}</div>
        </div>
      </div>
      <p class="muted">已登录，可直接发布文章、互动与私信。</p>
    `;
    els.logoutBtn.disabled = false;
    return;
  }

  els.authPanel.innerHTML = `
    <div class="auth-switch">
      <button class="chip ${authMode === 'login' ? 'chip-active' : ''}" data-auth-mode="login">登录</button>
      <button class="chip ${authMode === 'register' ? 'chip-active' : ''}" data-auth-mode="register">注册</button>
    </div>
    <form id="auth-form" class="stack-form">
      <input name="username" placeholder="用户名" required minlength="2" maxlength="24" />
      <input name="password" placeholder="密码（至少 6 位）" required type="password" minlength="6" />
      <button type="submit">${authMode === 'login' ? '登录并进入社区' : '注册并进入社区'}</button>
    </form>
  `;
  els.logoutBtn.disabled = true;
}

function renderFeedTabs() {
  const recommendActive = state.activeFeedTab === 'recommend' ? 'chip-active' : '';
  const followingActive = state.activeFeedTab === 'following' ? 'chip-active' : '';
  els.feedTabs.innerHTML = `
    <button class="chip ${recommendActive}" data-feed-tab="recommend">推荐</button>
    <button class="chip ${followingActive}" data-feed-tab="following">关注</button>
  `;
}

function renderFeedList() {
  const keyword = state.keyword.trim().toLowerCase();
  const list = activeFeed().filter((item) => {
    if (!keyword) {
      return true;
    }
    return (
      item.title.toLowerCase().includes(keyword) ||
      item.content.toLowerCase().includes(keyword) ||
      item.author.username.toLowerCase().includes(keyword)
    );
  });

  if (!state.user) {
    els.feedList.innerHTML = '<div class="empty">登录后查看推荐与关注流。</div>';
    return;
  }

  if (!list.length) {
    const text = state.activeFeedTab === 'following'
      ? '当前关注流为空，先去右侧“发现用户”关注一些人。'
      : '推荐流为空，先发布第一篇文章。';
    els.feedList.innerHTML = `<div class="empty">${text}</div>`;
    return;
  }

  els.feedList.innerHTML = list
    .map(
      (item) => `
        <article class="feed-card">
          <div class="feed-meta">
            <span class="author">${escapeHTML(item.author.username)}</span>
            <span>${timeText(item.created_at)}</span>
          </div>
          <h3>${escapeHTML(item.title)}</h3>
          <p>${escapeHTML(item.content.slice(0, 140))}${item.content.length > 140 ? '...' : ''}</p>
          <div class="feed-actions">
            <span>${item.like_count} 赞同</span>
            <span>${item.comment_count} 评论</span>
            <button data-open-detail="${item.id}">查看详情</button>
          </div>
        </article>
      `,
    )
    .join('');
}

function renderDetail() {
  if (!state.user) {
    els.detailPanel.innerHTML = '<div class="empty">登录后可查看文章详情与评论。</div>';
    return;
  }

  if (!state.articleDetail) {
    els.detailPanel.innerHTML = '<div class="empty">从信息流选择一篇文章开始互动。</div>';
    return;
  }

  const detail = state.articleDetail;
  const comments = detail.comments || [];

  els.detailPanel.innerHTML = `
    <div class="detail-head">
      <div>
        <h2>${escapeHTML(detail.article.title)}</h2>
        <div class="feed-meta">
          <span class="author">${escapeHTML(detail.author.username)}</span>
          <span>${timeText(detail.article.created_at)}</span>
        </div>
      </div>
      <button class="${detail.liked ? 'liked' : ''}" data-like-toggle="${detail.article.id}">
        ${detail.liked ? '已赞同' : '赞同'} (${detail.like_count})
      </button>
    </div>
    <div class="detail-content">${escapeHTML(detail.article.content)}</div>
    <div class="comments">
      <h4>评论 (${comments.length})</h4>
      <form id="comment-form" class="stack-form inline-form">
        <input name="content" placeholder="写下你的评论" required maxlength="280" />
        <button type="submit">发布</button>
      </form>
      <div class="comment-list">
        ${
          comments.length
            ? comments
                .map(
                  (comment) => `
              <div class="comment-item">
                <div class="feed-meta">
                  <span>用户 ${comment.user_id}</span>
                  <span>${timeText(comment.created_at)}</span>
                </div>
                <p>${escapeHTML(comment.content)}</p>
              </div>
            `,
                )
                .join('')
            : '<div class="muted">还没有评论，抢个沙发吧。</div>'
        }
      </div>
    </div>
  `;
}

function renderDiscoverUsers() {
  if (!state.user) {
    els.discoverList.innerHTML = '<div class="empty">登录后可发现用户并关注。</div>';
    return;
  }

  if (!state.discoverUsers.length) {
    els.discoverList.innerHTML = '<div class="empty">暂无可发现用户。</div>';
    return;
  }

  els.discoverList.innerHTML = state.discoverUsers
    .map(
      (user) => `
        <div class="user-row">
          <div>
            <div class="value">${escapeHTML(user.username)}</div>
            <div class="muted">UID ${user.id} ${user.is_mutual_follow ? '· 已互关' : ''}</div>
          </div>
          <button data-follow-id="${user.id}" ${user.is_following ? 'disabled' : ''}>
            ${user.is_following ? '已关注' : '关注'}
          </button>
        </div>
      `,
    )
    .join('');
}

function renderMutualUsers() {
  if (!state.user) {
    els.mutualList.innerHTML = '<div class="empty">登录后可查看互关联系人。</div>';
    return;
  }

  if (!state.mutualUsers.length) {
    els.mutualList.innerHTML = '<div class="empty">暂无互关用户，无法私信。</div>';
    return;
  }

  els.mutualList.innerHTML = state.mutualUsers
    .map((user) => {
      const activeClass = state.activePeer && state.activePeer.id === user.id ? 'row-active' : '';
      return `
        <button class="user-row user-select ${activeClass}" data-peer-id="${user.id}" data-peer-name="${escapeHTML(user.username)}">
          <span>${escapeHTML(user.username)}</span>
          <span class="muted">UID ${user.id}</span>
        </button>
      `;
    })
    .join('');
}

function renderChatPanel() {
  if (!state.user) {
    els.chatPanel.innerHTML = '<div class="empty">登录后使用私信。</div>';
    return;
  }

  if (!state.activePeer) {
    els.chatPanel.innerHTML = '<div class="empty">先选择一个互关联系人。</div>';
    return;
  }

  const rows = state.messages.length
    ? state.messages
        .map((msg) => {
          const mine = msg.from_user_id === state.user.id;
          return `
            <div class="chat-item ${mine ? 'mine' : ''}">
              <div class="chat-bubble">${escapeHTML(msg.content)}</div>
              <div class="muted">${mine ? '我' : state.activePeer.username} · ${timeText(msg.created_at)}</div>
            </div>
          `;
        })
        .join('')
    : '<div class="empty">还没有消息，打个招呼吧。</div>';

  els.chatPanel.innerHTML = `
    <div class="chat-head">与 ${escapeHTML(state.activePeer.username)} 的会话</div>
    <div class="chat-list">${rows}</div>
  `;
}

function renderAll() {
  els.apiBaseInput.value = state.apiBase;
  renderAuthPanel();
  renderFeedTabs();
  renderFeedList();
  renderDetail();
  renderDiscoverUsers();
  renderMutualUsers();
  renderChatPanel();
}

async function reloadFeeds() {
  const [recommend, following] = await Promise.all([recommendFeed(30), followingFeed()]);
  state.recommendFeed = recommend;
  state.followingFeed = following;
}

async function reloadSocial() {
  const [discover, mutual] = await Promise.all([discoverUsers(), mutualUsers()]);
  state.discoverUsers = discover;
  state.mutualUsers = mutual;

  if (state.activePeer) {
    const stillExists = mutual.find((u) => u.id === state.activePeer.id);
    if (!stillExists) {
      state.activePeer = null;
      state.messages = [];
      stopConversationPolling();
    }
  }
}

async function reloadConversation() {
  if (!state.activePeer) {
    state.messages = [];
    return;
  }
  state.messages = await conversation(state.activePeer.id);
}

async function loadDashboard() {
  await Promise.all([reloadFeeds(), reloadSocial()]);
  await reloadConversation();
}

async function openArticle(id) {
  state.articleDetail = await articleDetail(id);
  renderDetail();
}

async function onAuthSubmit(event) {
  event.preventDefault();
  const formData = new FormData(event.target);
  const username = String(formData.get('username') || '').trim();
  const password = String(formData.get('password') || '');

  if (!username || !password) {
    notify('请输入用户名和密码', 'warn');
    return;
  }

  try {
    const result = authMode === 'login' ? await login(username, password) : await register(username, password);
    setToken(result.token);
    setUser(result.user);
    state.activeFeedTab = 'recommend';
    notify(authMode === 'login' ? '登录成功' : '注册成功', 'ok');
    await loadDashboard();
    renderAll();
  } catch (err) {
    notify(err.message, 'error');
    renderAll();
  }
}

async function onPublish(event) {
  event.preventDefault();
  if (!requireLogin()) {
    return;
  }

  const formData = new FormData(event.target);
  const title = String(formData.get('title') || '').trim();
  const content = String(formData.get('content') || '').trim();

  if (!title || !content) {
    notify('标题和正文不能为空', 'warn');
    return;
  }

  try {
    await publishArticle(title, content);
    notify('发布成功，已更新推荐流', 'ok');
    event.target.reset();
    state.activeFeedTab = 'recommend';
    await reloadFeeds();
    renderAll();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onFeedClick(event) {
  const detailID = event.target.getAttribute('data-open-detail');
  if (!detailID) {
    return;
  }

  if (!requireLogin()) {
    return;
  }

  try {
    await openArticle(Number(detailID));
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onDetailClick(event) {
  const articleID = event.target.getAttribute('data-like-toggle');
  if (!articleID) {
    return;
  }

  if (!requireLogin() || !state.articleDetail) {
    return;
  }

  try {
    if (state.articleDetail.liked) {
      await unlikeArticle(Number(articleID));
    } else {
      await likeArticle(Number(articleID));
    }
    await openArticle(Number(articleID));
    await reloadFeeds();
    renderAll();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onCommentSubmit(event) {
  if (event.target.id !== 'comment-form') {
    return;
  }
  event.preventDefault();

  if (!requireLogin() || !state.articleDetail) {
    return;
  }

  const formData = new FormData(event.target);
  const content = String(formData.get('content') || '').trim();
  if (!content) {
    notify('评论内容不能为空', 'warn');
    return;
  }

  try {
    await addComment(state.articleDetail.article.id, content);
    event.target.reset();
    await openArticle(state.articleDetail.article.id);
    await reloadFeeds();
    renderAll();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onDiscoverClick(event) {
  const userID = event.target.getAttribute('data-follow-id');
  if (!userID) {
    return;
  }

  if (!requireLogin()) {
    return;
  }

  try {
    await followUser(Number(userID));
    notify('关注成功', 'ok');
    await Promise.all([reloadSocial(), reloadFeeds()]);
    renderAll();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onMutualClick(event) {
  const button = event.target.closest('[data-peer-id]');
  if (!button) {
    return;
  }

  const peerID = Number(button.getAttribute('data-peer-id'));
  const peerName = button.getAttribute('data-peer-name') || `用户 ${peerID}`;
  state.activePeer = { id: peerID, username: peerName };
  stopConversationPolling();

  try {
    await reloadConversation();
    renderAll();
    startConversationPolling();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onChatSubmit(event) {
  event.preventDefault();

  if (!requireLogin() || !state.activePeer) {
    return;
  }

  const formData = new FormData(event.target);
  const content = String(formData.get('content') || '').trim();
  if (!content) {
    notify('消息内容不能为空', 'warn');
    return;
  }

  try {
    await sendMessage(state.activePeer.id, content);
    event.target.reset();
    await reloadConversation();
    renderAll();
    startConversationPolling();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onRefresh() {
  if (!requireLogin()) {
    return;
  }

  try {
    await loadDashboard();
    if (state.articleDetail?.article?.id) {
      await openArticle(state.articleDetail.article.id);
    }
    renderAll();
    notify('数据已刷新', 'ok');
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onInitAuth() {
  if (!state.token) {
    renderAll();
    return;
  }

  try {
    const currentUser = await me();
    setUser(currentUser);
    await loadDashboard();
  } catch (err) {
    clearAuthState();
    notify(`登录状态已失效：${err.message}`, 'warn');
  }

  renderAll();
  startConversationPolling();
}

function bindEvents() {
  els.apiBaseSave.addEventListener('click', async () => {
    setApiBase(els.apiBaseInput.value);
    notify('API 地址已更新', 'ok');

    if (!state.token) {
      return;
    }

    try {
      await loadDashboard();
      renderAll();
      startConversationPolling();
    } catch (err) {
      notify(`加载失败：${err.message}`, 'error');
    }
  });

  els.authPanel.addEventListener('click', (event) => {
    const mode = event.target.getAttribute('data-auth-mode');
    if (!mode) {
      return;
    }
    authMode = mode;
    renderAuthPanel();
  });

  els.authPanel.addEventListener('submit', (event) => {
    if (event.target.id === 'auth-form') {
      onAuthSubmit(event);
    }
  });

  els.logoutBtn.addEventListener('click', () => {
    stopConversationPolling();
    clearAuthState();
    state.recommendFeed = [];
    state.followingFeed = [];
    state.discoverUsers = [];
    state.mutualUsers = [];
    state.activePeer = null;
    notify('已退出登录', 'ok');
    renderAll();
  });

  els.publishForm.addEventListener('submit', onPublish);

  els.feedTabs.addEventListener('click', (event) => {
    const tab = event.target.getAttribute('data-feed-tab');
    if (!tab) {
      return;
    }
    state.activeFeedTab = tab;
    renderFeedTabs();
    renderFeedList();
  });

  els.searchInput.addEventListener('input', (event) => {
    state.keyword = event.target.value || '';
    renderFeedList();
  });

  els.refreshBtn.addEventListener('click', onRefresh);
  els.feedList.addEventListener('click', onFeedClick);
  els.detailPanel.addEventListener('click', onDetailClick);
  els.detailPanel.addEventListener('submit', onCommentSubmit);
  els.discoverList.addEventListener('click', onDiscoverClick);
  els.mutualList.addEventListener('click', onMutualClick);
  els.chatForm.addEventListener('submit', onChatSubmit);

  document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
      stopConversationPolling();
      return;
    }
    startConversationPolling();
  });
}

function init() {
  bindEvents();
  onInitAuth();
}

init();
