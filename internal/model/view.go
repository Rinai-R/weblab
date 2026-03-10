package model

import "time"

type UserBrief struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserRelation struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	IsFollowing    bool   `json:"is_following"`
	IsMutualFollow bool   `json:"is_mutual_follow"`
}

type ArticleCard struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Author       UserBrief `json:"author"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type ArticleDetail struct {
	Article   Article   `json:"article"`
	Author    UserBrief `json:"author"`
	LikeCount int       `json:"like_count"`
	Liked     bool      `json:"liked"`
	Comments  []Comment `json:"comments"`
}

type QuestionCard struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      UserBrief `json:"author"`
	AnswerCount int       `json:"answer_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type AnswerView struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    UserBrief `json:"author"`
	VoteCount int       `json:"vote_count"`
	Voted     bool      `json:"voted"`
	CreatedAt time.Time `json:"created_at"`
}

type QuestionDetail struct {
	Question    Question     `json:"question"`
	Author      UserBrief    `json:"author"`
	AnswerCount int          `json:"answer_count"`
	Answers     []AnswerView `json:"answers"`
}
