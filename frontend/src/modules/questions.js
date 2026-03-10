import { request } from '../api/client.js';

export function askQuestion(title, description) {
  return request('/questions', {
    method: 'POST',
    data: { title, description },
  });
}

export function recommendQuestions(limit = 10, cursor = 0) {
  return request(`/questions/recommend?limit=${encodeURIComponent(limit)}&cursor=${encodeURIComponent(cursor)}`);
}

export function followingQuestions(limit = 10, cursor = 0) {
  return request(`/questions/following?limit=${encodeURIComponent(limit)}&cursor=${encodeURIComponent(cursor)}`);
}

export function questionDetail(id) {
  return request(`/questions/${id}`);
}

export function answerQuestion(questionID, content) {
  return request(`/questions/${questionID}/answers`, {
    method: 'POST',
    data: { content },
  });
}

export function questionAnswers(questionID, limit = 10, cursor = 0) {
  return request(
    `/questions/${questionID}/answers?limit=${encodeURIComponent(limit)}&cursor=${encodeURIComponent(cursor)}`,
  );
}

export function voteAnswer(answerID) {
  return request(`/questions/answers/${answerID}/vote`, { method: 'POST' });
}

export function unvoteAnswer(answerID) {
  return request(`/questions/answers/${answerID}/vote`, { method: 'DELETE' });
}
