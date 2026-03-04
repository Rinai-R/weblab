package mysql

import (
	"weblab/internal/model"

	"gorm.io/gorm"
)

type MessageDAO struct {
	db *gorm.DB
}

func NewMessageDAO(db *gorm.DB) *MessageDAO {
	return &MessageDAO{db: db}
}

func (d *MessageDAO) Create(message model.Message) (model.Message, error) {
	record := MessageRecord{
		FromUserID: message.FromUserID,
		ToUserID:   message.ToUserID,
		Content:    message.Content,
	}
	if err := d.db.Create(&record).Error; err != nil {
		return model.Message{}, normalizeErr(err)
	}

	return model.Message{
		ID:         record.ID,
		FromUserID: record.FromUserID,
		ToUserID:   record.ToUserID,
		Content:    record.Content,
		CreatedAt:  record.CreatedAt,
	}, nil
}

func (d *MessageDAO) ListConversation(userA, userB int64) ([]model.Message, error) {
	var records []MessageRecord
	err := d.db.Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userA, userB, userB, userA,
	).Order("created_at ASC").Order("id ASC").Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	messages := make([]model.Message, 0, len(records))
	for _, record := range records {
		messages = append(messages, model.Message{
			ID:         record.ID,
			FromUserID: record.FromUserID,
			ToUserID:   record.ToUserID,
			Content:    record.Content,
			CreatedAt:  record.CreatedAt,
		})
	}
	return messages, nil
}
