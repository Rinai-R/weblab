package mysql

import (
	"weblab/internal/model"

	"gorm.io/gorm"
)

type ArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) *ArticleDAO {
	return &ArticleDAO{db: db}
}

func (d *ArticleDAO) Create(article model.Article) (model.Article, error) {
	record := ArticleRecord{
		AuthorID: article.AuthorID,
		Title:    article.Title,
		Content:  article.Content,
	}
	if err := d.db.Create(&record).Error; err != nil {
		return model.Article{}, normalizeErr(err)
	}

	return model.Article{
		ID:        record.ID,
		AuthorID:  record.AuthorID,
		Title:     record.Title,
		Content:   record.Content,
		CreatedAt: record.CreatedAt,
	}, nil
}

func (d *ArticleDAO) GetByID(id int64) (model.Article, error) {
	var record ArticleRecord
	if err := d.db.Where("id = ?", id).First(&record).Error; err != nil {
		return model.Article{}, normalizeErr(err)
	}

	return model.Article{
		ID:        record.ID,
		AuthorID:  record.AuthorID,
		Title:     record.Title,
		Content:   record.Content,
		CreatedAt: record.CreatedAt,
	}, nil
}

func (d *ArticleDAO) ListByAuthorIDs(authorIDs []int64) ([]model.Article, error) {
	if len(authorIDs) == 0 {
		return []model.Article{}, nil
	}

	var records []ArticleRecord
	err := d.db.Where("author_id IN ?", authorIDs).
		Order("created_at DESC").
		Order("id DESC").
		Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	articles := make([]model.Article, 0, len(records))
	for _, record := range records {
		articles = append(articles, model.Article{
			ID:        record.ID,
			AuthorID:  record.AuthorID,
			Title:     record.Title,
			Content:   record.Content,
			CreatedAt: record.CreatedAt,
		})
	}
	return articles, nil
}

func (d *ArticleDAO) ListRecent(limit int) ([]model.Article, error) {
	if limit <= 0 {
		return []model.Article{}, nil
	}

	var records []ArticleRecord
	err := d.db.Order("created_at DESC").Order("id DESC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	articles := make([]model.Article, 0, len(records))
	for _, record := range records {
		articles = append(articles, model.Article{
			ID:        record.ID,
			AuthorID:  record.AuthorID,
			Title:     record.Title,
			Content:   record.Content,
			CreatedAt: record.CreatedAt,
		})
	}
	return articles, nil
}
