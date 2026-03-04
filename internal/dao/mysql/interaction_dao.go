package mysql

import (
	"time"
	"weblab/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractionDAO struct {
	db *gorm.DB
}

func NewInteractionDAO(db *gorm.DB) *InteractionDAO {
	return &InteractionDAO{db: db}
}

func (d *InteractionDAO) LikeArticle(userID, articleID int64) error {
	now := time.Now()
	record := ArticleLikeRecord{ArticleID: articleID, UserID: userID, CreatedAt: now}
	err := d.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "article_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"created_at": now,
		}),
	}).Create(&record).Error
	return normalizeErr(err)
}

func (d *InteractionDAO) UnlikeArticle(userID, articleID int64) error {
	err := d.db.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&ArticleLikeRecord{}).Error
	return normalizeErr(err)
}

func (d *InteractionDAO) IsArticleLiked(userID, articleID int64) (bool, error) {
	var count int64
	err := d.db.Model(&ArticleLikeRecord{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	if err != nil {
		return false, normalizeErr(err)
	}
	return count > 0, nil
}

func (d *InteractionDAO) CountLikes(articleID int64) (int, error) {
	var count int64
	err := d.db.Model(&ArticleLikeRecord{}).Where("article_id = ?", articleID).Count(&count).Error
	if err != nil {
		return 0, normalizeErr(err)
	}
	return int(count), nil
}

func (d *InteractionDAO) AddComment(comment model.Comment) (model.Comment, error) {
	record := CommentRecord{
		ArticleID: comment.ArticleID,
		UserID:    comment.UserID,
		Content:   comment.Content,
	}
	if err := d.db.Create(&record).Error; err != nil {
		return model.Comment{}, normalizeErr(err)
	}

	return model.Comment{
		ID:        record.ID,
		ArticleID: record.ArticleID,
		UserID:    record.UserID,
		Content:   record.Content,
		CreatedAt: record.CreatedAt,
	}, nil
}

func (d *InteractionDAO) ListComments(articleID int64) ([]model.Comment, error) {
	var records []CommentRecord
	err := d.db.Where("article_id = ?", articleID).
		Order("created_at ASC").
		Order("id ASC").
		Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	comments := make([]model.Comment, 0, len(records))
	for _, record := range records {
		comments = append(comments, model.Comment{
			ID:        record.ID,
			ArticleID: record.ArticleID,
			UserID:    record.UserID,
			Content:   record.Content,
			CreatedAt: record.CreatedAt,
		})
	}
	return comments, nil
}
