package service

import (
	"strings"
	"weblab/internal/dao"
	"weblab/internal/model"
)

type QuestionService struct {
	questionDAO dao.QuestionDAO
	userDAO     dao.UserDAO
}

func NewQuestionService(questionDAO dao.QuestionDAO, userDAO dao.UserDAO) *QuestionService {
	return &QuestionService{
		questionDAO: questionDAO,
		userDAO:     userDAO,
	}
}

func (s *QuestionService) Ask(userID int64, title, description string) (model.Question, error) {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if userID <= 0 || title == "" || description == "" {
		return model.Question{}, ErrInvalidArgument
	}

	if _, err := s.userDAO.GetByID(userID); err != nil {
		return model.Question{}, err
	}

	return s.questionDAO.CreateQuestion(model.Question{
		AuthorID:    userID,
		Title:       title,
		Description: description,
	})
}

func (s *QuestionService) Recommend(limit int, cursor int64) ([]model.QuestionCard, error) {
	if limit <= 0 {
		limit = 20
	}
	if cursor < 0 {
		return nil, ErrInvalidArgument
	}

	questions, err := s.questionDAO.ListQuestions(limit, cursor)
	if err != nil {
		return nil, err
	}
	return s.toCards(questions)
}

func (s *QuestionService) FollowFeed(userID int64, limit int, cursor int64) ([]model.QuestionCard, error) {
	if userID <= 0 || limit <= 0 || cursor < 0 {
		return nil, ErrInvalidArgument
	}

	if _, err := s.userDAO.GetByID(userID); err != nil {
		return nil, err
	}

	following, err := s.userDAO.ListFollowing(userID)
	if err != nil {
		return nil, err
	}
	if len(following) == 0 {
		return []model.QuestionCard{}, nil
	}

	questions, err := s.questionDAO.ListQuestionsByAuthorIDs(following, limit, cursor)
	if err != nil {
		return nil, err
	}
	return s.toCards(questions)
}

func (s *QuestionService) GetDetail(viewerID, questionID int64) (model.QuestionDetail, error) {
	if viewerID <= 0 || questionID <= 0 {
		return model.QuestionDetail{}, ErrInvalidArgument
	}

	if _, err := s.userDAO.GetByID(viewerID); err != nil {
		return model.QuestionDetail{}, err
	}

	question, err := s.questionDAO.GetQuestionByID(questionID)
	if err != nil {
		return model.QuestionDetail{}, err
	}
	author, err := s.userDAO.GetByID(question.AuthorID)
	if err != nil {
		return model.QuestionDetail{}, err
	}

	answerCount, err := s.questionDAO.CountAnswers(questionID)
	if err != nil {
		return model.QuestionDetail{}, err
	}
	answers, err := s.questionDAO.ListAnswers(questionID, 10, 0)
	if err != nil {
		return model.QuestionDetail{}, err
	}
	answerViews, err := s.toAnswerViews(viewerID, answers)
	if err != nil {
		return model.QuestionDetail{}, err
	}

	return model.QuestionDetail{
		Question: question,
		Author: model.UserBrief{
			ID:       author.ID,
			Username: author.Username,
		},
		AnswerCount: answerCount,
		Answers:     answerViews,
	}, nil
}

func (s *QuestionService) Answer(userID, questionID int64, content string) (model.AnswerView, error) {
	content = strings.TrimSpace(content)
	if userID <= 0 || questionID <= 0 || content == "" {
		return model.AnswerView{}, ErrInvalidArgument
	}

	if _, err := s.userDAO.GetByID(userID); err != nil {
		return model.AnswerView{}, err
	}
	if _, err := s.questionDAO.GetQuestionByID(questionID); err != nil {
		return model.AnswerView{}, err
	}

	answer, err := s.questionDAO.CreateAnswer(model.Answer{
		QuestionID: questionID,
		AuthorID:   userID,
		Content:    content,
	})
	if err != nil {
		return model.AnswerView{}, err
	}

	author, err := s.userDAO.GetByID(answer.AuthorID)
	if err != nil {
		return model.AnswerView{}, err
	}

	return model.AnswerView{
		ID:      answer.ID,
		Content: answer.Content,
		Author: model.UserBrief{
			ID:       author.ID,
			Username: author.Username,
		},
		VoteCount: 0,
		Voted:     false,
		CreatedAt: answer.CreatedAt,
	}, nil
}

func (s *QuestionService) ListAnswers(viewerID, questionID int64, limit int, cursor int64) ([]model.AnswerView, error) {
	if viewerID <= 0 || questionID <= 0 || limit <= 0 || cursor < 0 {
		return nil, ErrInvalidArgument
	}

	if _, err := s.userDAO.GetByID(viewerID); err != nil {
		return nil, err
	}
	if _, err := s.questionDAO.GetQuestionByID(questionID); err != nil {
		return nil, err
	}

	answers, err := s.questionDAO.ListAnswers(questionID, limit, cursor)
	if err != nil {
		return nil, err
	}
	return s.toAnswerViews(viewerID, answers)
}

func (s *QuestionService) VoteAnswer(userID, answerID int64) error {
	if userID <= 0 || answerID <= 0 {
		return ErrInvalidArgument
	}
	if _, err := s.userDAO.GetByID(userID); err != nil {
		return err
	}
	if _, err := s.questionDAO.GetAnswerByID(answerID); err != nil {
		return err
	}
	return s.questionDAO.VoteAnswer(userID, answerID)
}

func (s *QuestionService) UnvoteAnswer(userID, answerID int64) error {
	if userID <= 0 || answerID <= 0 {
		return ErrInvalidArgument
	}
	if _, err := s.userDAO.GetByID(userID); err != nil {
		return err
	}
	if _, err := s.questionDAO.GetAnswerByID(answerID); err != nil {
		return err
	}
	return s.questionDAO.UnvoteAnswer(userID, answerID)
}

func (s *QuestionService) toCards(questions []model.Question) ([]model.QuestionCard, error) {
	cards := make([]model.QuestionCard, 0, len(questions))
	for _, question := range questions {
		author, err := s.userDAO.GetByID(question.AuthorID)
		if err != nil {
			return nil, err
		}
		answerCount, err := s.questionDAO.CountAnswers(question.ID)
		if err != nil {
			return nil, err
		}

		cards = append(cards, model.QuestionCard{
			ID:          question.ID,
			Title:       question.Title,
			Description: question.Description,
			Author: model.UserBrief{
				ID:       author.ID,
				Username: author.Username,
			},
			AnswerCount: answerCount,
			CreatedAt:   question.CreatedAt,
		})
	}
	return cards, nil
}

func (s *QuestionService) toAnswerViews(viewerID int64, answers []model.Answer) ([]model.AnswerView, error) {
	views := make([]model.AnswerView, 0, len(answers))
	for _, answer := range answers {
		author, err := s.userDAO.GetByID(answer.AuthorID)
		if err != nil {
			return nil, err
		}
		voteCount, err := s.questionDAO.CountAnswerVotes(answer.ID)
		if err != nil {
			return nil, err
		}
		voted, err := s.questionDAO.IsAnswerVoted(viewerID, answer.ID)
		if err != nil {
			return nil, err
		}

		views = append(views, model.AnswerView{
			ID:      answer.ID,
			Content: answer.Content,
			Author: model.UserBrief{
				ID:       author.ID,
				Username: author.Username,
			},
			VoteCount: voteCount,
			Voted:     voted,
			CreatedAt: answer.CreatedAt,
		})
	}
	return views, nil
}
