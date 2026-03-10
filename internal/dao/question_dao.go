package dao

import "weblab/internal/model"

type QuestionDAO interface {
	CreateQuestion(question model.Question) (model.Question, error)
	GetQuestionByID(id int64) (model.Question, error)
	ListQuestions(limit int, cursor int64) ([]model.Question, error)
	ListQuestionsByAuthorIDs(authorIDs []int64, limit int, cursor int64) ([]model.Question, error)

	CreateAnswer(answer model.Answer) (model.Answer, error)
	GetAnswerByID(id int64) (model.Answer, error)
	ListAnswers(questionID int64, limit int, cursor int64) ([]model.Answer, error)
	CountAnswers(questionID int64) (int, error)

	VoteAnswer(userID, answerID int64) error
	UnvoteAnswer(userID, answerID int64) error
	IsAnswerVoted(userID, answerID int64) (bool, error)
	CountAnswerVotes(answerID int64) (int, error)
}
