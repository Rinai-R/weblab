package dao

import "weblab/internal/model"

type InteractionDAO interface {
	LikeArticle(userID, articleID int64) error
	UnlikeArticle(userID, articleID int64) error
	IsArticleLiked(userID, articleID int64) (bool, error)
	CountLikes(articleID int64) (int, error)
	AddComment(comment model.Comment) (model.Comment, error)
	ListComments(articleID int64) ([]model.Comment, error)
}
