package mysql

import "gorm.io/gorm"

func EnsureSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&UserRecord{},
		&ArticleRecord{},
		&QuestionRecord{},
		&FollowRecord{},
		&ArticleLikeRecord{},
		&CommentRecord{},
		&AnswerRecord{},
		&AnswerVoteRecord{},
		&MessageRecord{},
	)
}
