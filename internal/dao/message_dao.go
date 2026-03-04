package dao

import "weblab/internal/model"

type MessageDAO interface {
	Create(message model.Message) (model.Message, error)
	ListConversation(userA, userB int64) ([]model.Message, error)
}
