package mysql

import "time"

type UserRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"size:64;not null;uniqueIndex"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UserRecord) TableName() string { return "users" }

type ArticleRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	AuthorID  int64     `gorm:"column:author_id;not null;index:idx_articles_author_created,priority:1"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_articles_author_created,priority:2"`
}

func (ArticleRecord) TableName() string { return "articles" }

type FollowRecord struct {
	FollowerID int64     `gorm:"column:follower_id;primaryKey;autoIncrement:false"`
	FolloweeID int64     `gorm:"column:followee_id;primaryKey;autoIncrement:false;index:idx_follows_followee"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (FollowRecord) TableName() string { return "follows" }

type ArticleLikeRecord struct {
	ArticleID int64     `gorm:"column:article_id;primaryKey;autoIncrement:false"`
	UserID    int64     `gorm:"column:user_id;primaryKey;autoIncrement:false;index:idx_article_likes_user"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ArticleLikeRecord) TableName() string { return "article_likes" }

type CommentRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	ArticleID int64     `gorm:"column:article_id;not null;index:idx_comments_article_created,priority:1"`
	UserID    int64     `gorm:"column:user_id;not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_comments_article_created,priority:2"`
}

func (CommentRecord) TableName() string { return "comments" }

type MessageRecord struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	FromUserID int64     `gorm:"column:from_user_id;not null;index:idx_messages_pair_created,priority:1;index:idx_messages_reverse_pair_created,priority:2"`
	ToUserID   int64     `gorm:"column:to_user_id;not null;index:idx_messages_pair_created,priority:2;index:idx_messages_reverse_pair_created,priority:1"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_messages_pair_created,priority:3;index:idx_messages_reverse_pair_created,priority:3"`
}

func (MessageRecord) TableName() string { return "messages" }
