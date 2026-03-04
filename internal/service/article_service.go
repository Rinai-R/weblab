package service

import (
	"strings"
	"weblab/internal/dao"
	"weblab/internal/model"
)

type ArticleService struct {
	articleDAO     dao.ArticleDAO
	userDAO        dao.UserDAO
	interactionDAO dao.InteractionDAO
}

func NewArticleService(articleDAO dao.ArticleDAO, userDAO dao.UserDAO, interactionDAO dao.InteractionDAO) *ArticleService {
	return &ArticleService{articleDAO: articleDAO, userDAO: userDAO, interactionDAO: interactionDAO}
}

func (s *ArticleService) Publish(authorID int64, title, content string) (model.Article, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" || content == "" {
		return model.Article{}, ErrInvalidArgument
	}

	if _, err := s.userDAO.GetByID(authorID); err != nil {
		return model.Article{}, err
	}

	return s.articleDAO.Create(model.Article{AuthorID: authorID, Title: title, Content: content})
}

func (s *ArticleService) GetDetail(userID, articleID int64) (model.ArticleDetail, error) {
	if userID <= 0 || articleID <= 0 {
		return model.ArticleDetail{}, ErrInvalidArgument
	}
	if _, err := s.userDAO.GetByID(userID); err != nil {
		return model.ArticleDetail{}, err
	}

	article, err := s.articleDAO.GetByID(articleID)
	if err != nil {
		return model.ArticleDetail{}, err
	}

	author, err := s.userDAO.GetByID(article.AuthorID)
	if err != nil {
		return model.ArticleDetail{}, err
	}

	likes, err := s.interactionDAO.CountLikes(articleID)
	if err != nil {
		return model.ArticleDetail{}, err
	}

	liked, err := s.interactionDAO.IsArticleLiked(userID, articleID)
	if err != nil {
		return model.ArticleDetail{}, err
	}

	comments, err := s.interactionDAO.ListComments(articleID)
	if err != nil {
		return model.ArticleDetail{}, err
	}

	return model.ArticleDetail{
		Article: article,
		Author: model.UserBrief{
			ID:       author.ID,
			Username: author.Username,
		},
		LikeCount: likes,
		Liked:     liked,
		Comments:  comments,
	}, nil
}

func (s *ArticleService) Feed(userID int64) ([]model.ArticleCard, error) {
	if userID <= 0 {
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
		return []model.ArticleCard{}, nil
	}
	articles, err := s.articleDAO.ListByAuthorIDs(following)
	if err != nil {
		return nil, err
	}
	return s.toCards(articles)
}

func (s *ArticleService) Recommend(limit int) ([]model.ArticleCard, error) {
	if limit <= 0 {
		limit = 20
	}

	articles, err := s.articleDAO.ListRecent(limit)
	if err != nil {
		return nil, err
	}
	return s.toCards(articles)
}

func (s *ArticleService) toCards(articles []model.Article) ([]model.ArticleCard, error) {
	cards := make([]model.ArticleCard, 0, len(articles))
	for _, article := range articles {
		author, err := s.userDAO.GetByID(article.AuthorID)
		if err != nil {
			return nil, err
		}

		likes, err := s.interactionDAO.CountLikes(article.ID)
		if err != nil {
			return nil, err
		}
		comments, err := s.interactionDAO.ListComments(article.ID)
		if err != nil {
			return nil, err
		}

		cards = append(cards, model.ArticleCard{
			ID:      article.ID,
			Title:   article.Title,
			Content: article.Content,
			Author: model.UserBrief{
				ID:       author.ID,
				Username: author.Username,
			},
			LikeCount:    likes,
			CommentCount: len(comments),
			CreatedAt:    article.CreatedAt,
		})
	}
	return cards, nil
}
