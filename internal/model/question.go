package model

import "time"

type Question struct {
	ID          int64     `json:"id"`
	AuthorID    int64     `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Answer struct {
	ID         int64     `json:"id"`
	QuestionID int64     `json:"question_id"`
	AuthorID   int64     `json:"author_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
