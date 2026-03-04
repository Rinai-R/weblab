package service

import (
	"strings"
	"weblab/internal/dao"
	"weblab/internal/model"
)

type InteractionService struct {
	interactionDAO dao.InteractionDAO
	articleDAO     dao.ArticleDAO
}

func NewInteractionService(interactionDAO dao.InteractionDAO, articleDAO dao.ArticleDAO) *InteractionService {
	return &InteractionService{interactionDAO: interactionDAO, articleDAO: articleDAO}
}

func (s *InteractionService) LikeArticle(userID, articleID int64) error {
	if userID <= 0 || articleID <= 0 {
		return ErrInvalidArgument
	}
	if _, err := s.articleDAO.GetByID(articleID); err != nil {
		return err
	}
	return s.interactionDAO.LikeArticle(userID, articleID)
}

func (s *InteractionService) UnlikeArticle(userID, articleID int64) error {
	if userID <= 0 || articleID <= 0 {
		return ErrInvalidArgument
	}
	if _, err := s.articleDAO.GetByID(articleID); err != nil {
		return err
	}
	return s.interactionDAO.UnlikeArticle(userID, articleID)
}

func (s *InteractionService) AddComment(userID, articleID int64, content string) (model.Comment, error) {
	content = strings.TrimSpace(content)
	if userID <= 0 || articleID <= 0 || content == "" {
		return model.Comment{}, ErrInvalidArgument
	}
	if _, err := s.articleDAO.GetByID(articleID); err != nil {
		return model.Comment{}, err
	}
	return s.interactionDAO.AddComment(model.Comment{ArticleID: articleID, UserID: userID, Content: content})
}

func (s *InteractionService) ListComments(articleID int64) ([]model.Comment, error) {
	if articleID <= 0 {
		return nil, ErrInvalidArgument
	}
	if _, err := s.articleDAO.GetByID(articleID); err != nil {
		return nil, err
	}
	return s.interactionDAO.ListComments(articleID)
}
