package dao

import "weblab/internal/model"

type ArticleDAO interface {
	Create(article model.Article) (model.Article, error)
	GetByID(id int64) (model.Article, error)
	ListByAuthorIDs(authorIDs []int64) ([]model.Article, error)
	ListRecent(limit int) ([]model.Article, error)
}
