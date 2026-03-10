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
  answerQuestion,
  askQuestion,
  followingQuestions,
  questionAnswers,
  questionDetail,
  recommendQuestions,
  unvoteAnswer,
  voteAnswer,
} from './modules/questions.js';
import {
  conversation,
  discoverUsers,
  followUser,
  mutualUsers,
  sendMessage,
} from './modules/social.js';

const FEED_PAGE_SIZE = 8;
const ANSWER_PAGE_SIZE = 6;
const CHAT_POLL_INTERVAL_MS = 2000;

const STREAM_LABELS = {
  articles_recommend: '推荐 · 文章',
  articles_following: '关注 · 文章',
  questions_recommend: '推荐 · 问答',
  questions_following: '关注 · 问答',
};

const STREAM_KEYS = Object.keys(STREAM_LABELS);

const els = {
  notice: document.getElementById('notice'),
  apiBaseInput: document.getElementById('api-base-input'),
  apiBaseSave: document.getElementById('api-base-save'),
  logoutBtn: document.getElementById('logout-btn'),
  authPanel: document.getElementById('auth-panel'),
  publishForm: document.getElementById('publish-form'),
  askForm: document.getElementById('ask-form'),
  streamNav: document.getElementById('stream-nav'),
  feedSearch: document.getElementById('feed-search-input'),
  refreshBtn: document.getElementById('stream-refresh-btn'),
  feedList: document.getElementById('feed-list'),
  loadMoreBtn: document.getElementById('feed-load-more-btn'),
  detailPanel: document.getElementById('detail-panel'),
  discoverList: document.getElementById('discover-list'),
  mutualList: document.getElementById('mutual-list'),
  chatPanel: document.getElementById('chat-panel'),
  chatForm: document.getElementById('chat-form'),
};

let authMode = 'login';
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

function streamState(streamKey = state.activeStream) {
  return state.streams[streamKey];
}

function currentItems() {
  return streamState().items || [];
}

function notify(message, type = 'info') {
  els.notice.textContent = message;
  els.notice.className = `notice notice-${type}`;

  window.clearTimeout(notify.timer);
  notify.timer = window.setTimeout(() => {
    els.notice.className = 'notice';
    els.notice.textContent = '';
  }, 3000);
}

function requireLogin() {
  if (state.token && state.user) {
    return true;
  }
  notify('请先登录后操作', 'warn');
  return false;
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
  return state.messages[state.messages.length - 1].id !== nextMessages[nextMessages.length - 1].id;
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
      const next = await conversation(peerID);
      if (state.activePeer && state.activePeer.id === peerID && hasConversationChanged(next)) {
        state.messages = next;
        renderChatPanel();
      }
    } catch {
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

function renderAuthPanel() {
  if (state.user) {
    els.authPanel.innerHTML = `
      <div class="user-head">
        <div class="avatar">${escapeHTML(state.user.username.slice(0, 1).toUpperCase())}</div>
        <div>
          <div><strong>${escapeHTML(state.user.username)}</strong></div>
          <div class="muted">UID ${state.user.id}</div>
        </div>
      </div>
      <p class="muted">支持发文、提问、回答、关注和私信。</p>
    `;
    els.logoutBtn.disabled = false;
    return;
  }

  els.authPanel.innerHTML = `
    <div class="auth-switch">
      <button data-auth-mode="login" class="${authMode === 'login' ? 'active' : ''}">登录</button>
      <button data-auth-mode="register" class="${authMode === 'register' ? 'active' : ''}">注册</button>
    </div>
    <form id="auth-form" class="stack-form">
      <input name="username" required maxlength="24" placeholder="用户名" />
      <input name="password" required minlength="6" type="password" placeholder="密码（最少 6 位）" />
      <button type="submit">${authMode === 'login' ? '登录' : '注册'}</button>
    </form>
  `;
  els.logoutBtn.disabled = true;
}

function renderStreamNav() {
  els.streamNav.innerHTML = STREAM_KEYS.map((key) => {
    const activeClass = state.activeStream === key ? 'active' : '';
    return `<button class="tab-btn ${activeClass}" data-stream="${key}">${STREAM_LABELS[key]}</button>`;
  }).join('');
}

function filteredItems() {
  const keyword = state.keyword.trim().toLowerCase();
  if (!keyword) {
    return currentItems();
  }

  return currentItems().filter((item) => {
    const title = String(item.title || '').toLowerCase();
    const content = String(item.content || item.description || '').toLowerCase();
    const username = String(item.author?.username || '').toLowerCase();
    return title.includes(keyword) || content.includes(keyword) || username.includes(keyword);
  });
}

function renderFeedList() {
  if (!state.user) {
    els.feedList.innerHTML = '<div class="empty">登录后查看推荐流、关注流和问答流。</div>';
    return;
  }

  const items = filteredItems();
  if (!items.length) {
    els.feedList.innerHTML = '<div class="empty">暂无内容，试试切换分栏或发布新内容。</div>';
    return;
  }

  if (state.activeStream.startsWith('articles_')) {
    els.feedList.innerHTML = items.map((item) => `
      <article class="feed-card">
        <div class="feed-meta">
          <span class="author">${escapeHTML(item.author.username)}</span>
          <span>${timeText(item.created_at)}</span>
        </div>
        <h4>${escapeHTML(item.title)}</h4>
        <p>${escapeHTML(item.content.slice(0, 180))}${item.content.length > 180 ? '...' : ''}</p>
        <div class="feed-actions">
          <span>${item.like_count} 赞同</span>
          <span>${item.comment_count} 评论</span>
          <button data-open-article="${item.id}">查看详情</button>
        </div>
      </article>
    `).join('');
    return;
  }

  els.feedList.innerHTML = items.map((item) => `
    <article class="feed-card">
      <div class="feed-meta">
        <span class="author">${escapeHTML(item.author.username)}</span>
        <span>${timeText(item.created_at)}</span>
      </div>
      <h4>${escapeHTML(item.title)}</h4>
      <p>${escapeHTML(item.description.slice(0, 180))}${item.description.length > 180 ? '...' : ''}</p>
      <div class="feed-actions">
        <span>${item.answer_count} 回答</span>
        <button data-open-question="${item.id}">查看问答</button>
      </div>
    </article>
  `).join('');
}

function renderLoadMoreButton() {
  const stream = streamState();
  if (!state.user) {
    els.loadMoreBtn.style.display = 'none';
    return;
  }

  els.loadMoreBtn.style.display = stream.hasMore ? 'block' : 'none';
}

function renderDetailPanel() {
  if (!state.user) {
    els.detailPanel.innerHTML = '<div class="empty">登录后可查看详情。</div>';
    return;
  }

  if (state.articleDetail) {
    const detail = state.articleDetail;
    const comments = detail.comments || [];

    els.detailPanel.innerHTML = `
      <div class="detail-head">
        <div>
          <h2 class="detail-title">${escapeHTML(detail.article.title)}</h2>
          <div class="feed-meta">
            <span class="author">${escapeHTML(detail.author.username)}</span>
            <span>${timeText(detail.article.created_at)}</span>
          </div>
        </div>
        <button data-like-article="${detail.article.id}">${detail.liked ? '已赞同' : '赞同'} (${detail.like_count})</button>
      </div>
      <div class="detail-content">${escapeHTML(detail.article.content)}</div>
      <div class="detail-sub">
        <h3>评论 (${comments.length})</h3>
        <form id="comment-form" class="stack-form inline-form">
          <input name="content" maxlength="280" required placeholder="写下你的评论" />
          <button type="submit">发布</button>
        </form>
        ${comments.length ? comments.map((comment) => `
          <div class="comment-item">
            <div class="feed-meta">
              <span>用户 ${comment.user_id}</span>
              <span>${timeText(comment.created_at)}</span>
            </div>
            <p>${escapeHTML(comment.content)}</p>
          </div>
        `).join('') : '<div class="empty">暂无评论，来抢沙发。</div>'}
      </div>
    `;
    return;
  }

  if (state.questionDetail) {
    const detail = state.questionDetail;
    const answers = detail.answers || [];

    els.detailPanel.innerHTML = `
      <div class="detail-head">
        <div>
          <h2 class="detail-title">${escapeHTML(detail.question.title)}</h2>
          <div class="feed-meta">
            <span class="author">${escapeHTML(detail.author.username)}</span>
            <span>${timeText(detail.question.created_at)}</span>
            <span>${detail.answer_count} 回答</span>
          </div>
        </div>
      </div>
      <div class="detail-content">${escapeHTML(detail.question.description)}</div>
      <div class="detail-sub">
        <h3>写回答</h3>
        <form id="answer-form" class="stack-form">
          <textarea name="content" required maxlength="5000" placeholder="分享你的看法与经验"></textarea>
          <button type="submit">发布回答</button>
        </form>
      </div>
      <div class="detail-sub">
        <h3>全部回答</h3>
        ${answers.length ? answers.map((answer) => `
          <div class="answer-item">
            <div class="feed-meta">
              <span class="author">${escapeHTML(answer.author.username)}</span>
              <span>${timeText(answer.created_at)}</span>
            </div>
            <p>${escapeHTML(answer.content)}</p>
            <div class="feed-actions">
              <span>${answer.vote_count} 赞同</span>
              <button data-vote-answer="${answer.id}" data-voted="${answer.voted ? '1' : '0'}">
                ${answer.voted ? '取消赞同' : '赞同'}
              </button>
            </div>
          </div>
        `).join('') : '<div class="empty">还没有回答，欢迎第一个作答。</div>'}
        ${state.answerHasMore ? '<button data-load-more-answers="1" class="load-more">加载更多回答</button>' : ''}
      </div>
    `;
    return;
  }

  els.detailPanel.innerHTML = '<div class="empty">点击信息流中的内容查看详情。</div>';
}

function renderDiscoverUsers() {
  if (!state.user) {
    els.discoverList.innerHTML = '<div class="empty">登录后可发现用户。</div>';
    return;
  }

  if (!state.discoverUsers.length) {
    els.discoverList.innerHTML = '<div class="empty">暂无推荐用户。</div>';
    return;
  }

  els.discoverList.innerHTML = state.discoverUsers.map((user) => `
    <div class="user-row">
      <div>
        <div><strong>${escapeHTML(user.username)}</strong></div>
        <div class="muted">UID ${user.id}${user.is_mutual_follow ? ' · 已互关' : ''}</div>
      </div>
      <button data-follow-id="${user.id}" ${user.is_following ? 'disabled' : ''}>${user.is_following ? '已关注' : '关注'}</button>
    </div>
  `).join('');
}

function renderMutualUsers() {
  if (!state.user) {
    els.mutualList.innerHTML = '<div class="empty">登录后可查看互关。</div>';
    return;
  }

  if (!state.mutualUsers.length) {
    els.mutualList.innerHTML = '<div class="empty">暂无互关用户。</div>';
    return;
  }

  els.mutualList.innerHTML = state.mutualUsers.map((user) => {
    const activeClass = state.activePeer?.id === user.id ? 'active' : '';
    return `
      <button class="user-select ${activeClass}" data-peer-id="${user.id}" data-peer-name="${escapeHTML(user.username)}">
        ${escapeHTML(user.username)} <span class="muted">UID ${user.id}</span>
      </button>
    `;
  }).join('');
}

function renderChatPanel() {
  if (!state.user) {
    els.chatPanel.innerHTML = '<div class="empty">登录后可聊天。</div>';
    return;
  }

  if (!state.activePeer) {
    els.chatPanel.innerHTML = '<div class="empty">选择互关用户开始聊天。</div>';
    return;
  }

  const rows = state.messages.length
    ? state.messages.map((msg) => {
      const mine = msg.from_user_id === state.user.id;
      return `
        <div class="chat-item ${mine ? 'mine' : ''}">
          <div class="chat-bubble">${escapeHTML(msg.content)}</div>
          <div class="muted">${mine ? '我' : state.activePeer.username} · ${timeText(msg.created_at)}</div>
        </div>
      `;
    }).join('')
    : '<div class="empty">会话暂无消息。</div>';

  els.chatPanel.innerHTML = `<div class="chat-list">${rows}</div>`;
}

function renderAll() {
  els.apiBaseInput.value = state.apiBase;
  renderAuthPanel();
  renderStreamNav();
  renderFeedList();
  renderLoadMoreButton();
  renderDetailPanel();
  renderDiscoverUsers();
  renderMutualUsers();
  renderChatPanel();
}

async function fetchStreamData(streamKey, cursor) {
  switch (streamKey) {
    case 'articles_recommend':
      return recommendFeed(FEED_PAGE_SIZE, cursor);
    case 'articles_following':
      return followingFeed(FEED_PAGE_SIZE, cursor);
    case 'questions_recommend':
      return recommendQuestions(FEED_PAGE_SIZE, cursor);
    case 'questions_following':
      return followingQuestions(FEED_PAGE_SIZE, cursor);
    default:
      throw new Error('unknown stream');
  }
}

async function loadStream(streamKey, { append = false } = {}) {
  if (!state.user) {
    return;
  }

  const stream = streamState(streamKey);
  const cursor = append ? stream.cursor : 0;
  const payload = await fetchStreamData(streamKey, cursor);
  const items = payload.items || [];
  const nextCursor = payload.next_cursor || 0;

  stream.items = append ? stream.items.concat(items) : items;
  stream.cursor = nextCursor;
  stream.hasMore = items.length === FEED_PAGE_SIZE && nextCursor > 0;
}

async function loadMoreActiveStream() {
  if (!requireLogin()) {
    return;
  }

  const stream = streamState();
  if (!stream.hasMore) {
    return;
  }

  try {
    await loadStream(state.activeStream, { append: true });
    renderFeedList();
    renderLoadMoreButton();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function openArticleDetail(id) {
  state.articleDetail = await articleDetail(id);
  state.questionDetail = null;
}

async function openQuestionDetail(id) {
  const detail = await questionDetail(id);
  state.questionDetail = detail;
  state.articleDetail = null;

  const answers = detail.answers || [];
  state.answerCursor = answers.length ? answers[answers.length - 1].id : 0;
  state.answerHasMore = answers.length < detail.answer_count;
}

async function loadMoreAnswers() {
  if (!state.questionDetail || !state.answerHasMore) {
    return;
  }

  const questionID = state.questionDetail.question.id;
  const payload = await questionAnswers(questionID, ANSWER_PAGE_SIZE, state.answerCursor);
  const items = payload.items || [];
  if (!items.length) {
    state.answerHasMore = false;
    return;
  }

  const existing = new Set((state.questionDetail.answers || []).map((x) => x.id));
  const deduped = items.filter((x) => !existing.has(x.id));
  state.questionDetail.answers = (state.questionDetail.answers || []).concat(deduped);
  state.answerCursor = payload.next_cursor || items[items.length - 1].id;
  state.answerHasMore = state.questionDetail.answers.length < state.questionDetail.answer_count;
}

async function loadSocial() {
  if (!state.user) {
    return;
  }

  const [discover, mutual] = await Promise.all([discoverUsers(), mutualUsers()]);
  state.discoverUsers = discover;
  state.mutualUsers = mutual;

  if (state.activePeer && !mutual.find((x) => x.id === state.activePeer.id)) {
    state.activePeer = null;
    state.messages = [];
    stopConversationPolling();
  }
}

async function reloadConversation() {
  if (!state.activePeer) {
    state.messages = [];
    return;
  }
  state.messages = await conversation(state.activePeer.id);
}

async function onAuthSubmit(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const username = String(formData.get('username') || '').trim();
  const password = String(formData.get('password') || '');

  if (!username || !password) {
    notify('用户名和密码不能为空', 'warn');
    return;
  }

  try {
    const result = authMode === 'login' ? await login(username, password) : await register(username, password);
    setToken(result.token);
    setUser(result.user);

    await Promise.all([
      loadStream(state.activeStream, { append: false }),
      loadSocial(),
    ]);

    state.articleDetail = null;
    state.questionDetail = null;

    renderAll();
    notify(authMode === 'login' ? '登录成功' : '注册成功', 'ok');
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onPublishArticle(event) {
  event.preventDefault();
  if (!requireLogin()) {
    return;
  }

  const formData = new FormData(event.target);
  const title = String(formData.get('title') || '').trim();
  const content = String(formData.get('content') || '').trim();

  if (!title || !content) {
    notify('标题和内容不能为空', 'warn');
    return;
  }

  try {
    await publishArticle(title, content);
    event.target.reset();
    state.activeStream = 'articles_recommend';
    await loadStream('articles_recommend');
    renderAll();
    notify('文章已发布', 'ok');
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onAskQuestion(event) {
  event.preventDefault();
  if (!requireLogin()) {
    return;
  }

  const formData = new FormData(event.target);
  const title = String(formData.get('title') || '').trim();
  const description = String(formData.get('description') || '').trim();

  if (!title || !description) {
    notify('问题标题和描述不能为空', 'warn');
    return;
  }

  try {
    await askQuestion(title, description);
    event.target.reset();
    state.activeStream = 'questions_recommend';
    await loadStream('questions_recommend');
    renderAll();
    notify('问题已发布', 'ok');
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onFeedClick(event) {
  const articleID = event.target.getAttribute('data-open-article');
  const questionID = event.target.getAttribute('data-open-question');

  if (!articleID && !questionID) {
    return;
  }
  if (!requireLogin()) {
    return;
  }

  try {
    if (articleID) {
      await openArticleDetail(Number(articleID));
    }
    if (questionID) {
      await openQuestionDetail(Number(questionID));
    }
    renderDetailPanel();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onDetailClick(event) {
  const likeArticleID = event.target.getAttribute('data-like-article');
  const voteAnswerID = event.target.getAttribute('data-vote-answer');
  const loadMoreAnswersFlag = event.target.getAttribute('data-load-more-answers');

  if (!likeArticleID && !voteAnswerID && !loadMoreAnswersFlag) {
    return;
  }
  if (!requireLogin()) {
    return;
  }

  try {
    if (likeArticleID && state.articleDetail) {
      if (state.articleDetail.liked) {
        await unlikeArticle(Number(likeArticleID));
      } else {
        await likeArticle(Number(likeArticleID));
      }
      await openArticleDetail(Number(likeArticleID));
      await loadStream(state.activeStream, { append: false });
      renderAll();
      return;
    }

    if (voteAnswerID && state.questionDetail) {
      const voted = event.target.getAttribute('data-voted') === '1';
      if (voted) {
        await unvoteAnswer(Number(voteAnswerID));
      } else {
        await voteAnswer(Number(voteAnswerID));
      }
      await openQuestionDetail(state.questionDetail.question.id);
      await loadStream(state.activeStream, { append: false });
      renderAll();
      return;
    }

    if (loadMoreAnswersFlag && state.questionDetail) {
      await loadMoreAnswers();
      renderDetailPanel();
    }
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onDetailSubmit(event) {
  if (event.target.id !== 'comment-form' && event.target.id !== 'answer-form') {
    return;
  }
  event.preventDefault();

  if (!requireLogin()) {
    return;
  }

  try {
    if (event.target.id === 'comment-form' && state.articleDetail) {
      const content = String(new FormData(event.target).get('content') || '').trim();
      if (!content) {
        notify('评论不能为空', 'warn');
        return;
      }

      await addComment(state.articleDetail.article.id, content);
      event.target.reset();
      await openArticleDetail(state.articleDetail.article.id);
      await loadStream(state.activeStream, { append: false });
      renderAll();
      return;
    }

    if (event.target.id === 'answer-form' && state.questionDetail) {
      const content = String(new FormData(event.target).get('content') || '').trim();
      if (!content) {
        notify('回答不能为空', 'warn');
        return;
      }

      await answerQuestion(state.questionDetail.question.id, content);
      event.target.reset();
      await openQuestionDetail(state.questionDetail.question.id);
      await loadStream(state.activeStream, { append: false });
      renderAll();
    }
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onFollowClick(event) {
  const followID = event.target.getAttribute('data-follow-id');
  if (!followID) {
    return;
  }
  if (!requireLogin()) {
    return;
  }

  try {
    await followUser(Number(followID));
    await Promise.all([loadSocial(), loadStream(state.activeStream, { append: false })]);
    renderAll();
    notify('关注成功', 'ok');
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onSelectPeer(event) {
  const button = event.target.closest('[data-peer-id]');
  if (!button) {
    return;
  }

  const peerID = Number(button.getAttribute('data-peer-id'));
  const peerName = button.getAttribute('data-peer-name') || `用户 ${peerID}`;
  state.activePeer = { id: peerID, username: peerName };

  try {
    await reloadConversation();
    renderMutualUsers();
    renderChatPanel();
    startConversationPolling();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onSendMessage(event) {
  event.preventDefault();
  if (!requireLogin() || !state.activePeer) {
    return;
  }

  const content = String(new FormData(event.target).get('content') || '').trim();
  if (!content) {
    notify('消息不能为空', 'warn');
    return;
  }

  try {
    await sendMessage(state.activePeer.id, content);
    event.target.reset();
    await reloadConversation();
    renderChatPanel();
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
    await Promise.all([
      loadStream(state.activeStream, { append: false }),
      loadSocial(),
    ]);

    if (state.articleDetail?.article?.id) {
      await openArticleDetail(state.articleDetail.article.id);
    }
    if (state.questionDetail?.question?.id) {
      await openQuestionDetail(state.questionDetail.question.id);
    }

    renderAll();
    notify('数据已刷新', 'ok');
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function onSwitchStream(streamKey) {
  if (!STREAM_LABELS[streamKey]) {
    return;
  }
  if (!requireLogin()) {
    return;
  }

  state.activeStream = streamKey;
  state.articleDetail = null;
  state.questionDetail = null;

  try {
    await loadStream(streamKey, { append: false });
    renderAll();
  } catch (err) {
    notify(err.message, 'error');
  }
}

async function initAuth() {
  if (!state.token) {
    renderAll();
    return;
  }

  try {
    const currentUser = await me();
    setUser(currentUser);
    await Promise.all([
      loadStream(state.activeStream, { append: false }),
      loadSocial(),
    ]);
  } catch (err) {
    clearAuthState();
    notify(`登录状态已失效：${err.message}`, 'warn');
  }

  renderAll();
}

function bindEvents() {
  els.apiBaseSave.addEventListener('click', async () => {
    setApiBase(els.apiBaseInput.value);
    notify('API 地址已更新', 'ok');

    if (!state.user) {
      return;
    }

    try {
      await Promise.all([
        loadStream(state.activeStream, { append: false }),
        loadSocial(),
      ]);
      renderAll();
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
    renderAll();
    notify('已退出登录', 'ok');
  });

  els.publishForm.addEventListener('submit', onPublishArticle);
  els.askForm.addEventListener('submit', onAskQuestion);

  els.streamNav.addEventListener('click', (event) => {
    const streamKey = event.target.getAttribute('data-stream');
    if (!streamKey) {
      return;
    }
    onSwitchStream(streamKey);
  });

  els.feedSearch.addEventListener('input', (event) => {
    state.keyword = event.target.value || '';
    renderFeedList();
  });

  els.refreshBtn.addEventListener('click', onRefresh);
  els.feedList.addEventListener('click', onFeedClick);
  els.loadMoreBtn.addEventListener('click', loadMoreActiveStream);

  els.detailPanel.addEventListener('click', onDetailClick);
  els.detailPanel.addEventListener('submit', onDetailSubmit);

  els.discoverList.addEventListener('click', onFollowClick);
  els.mutualList.addEventListener('click', onSelectPeer);
  els.chatForm.addEventListener('submit', onSendMessage);

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
  initAuth();
}

init();
