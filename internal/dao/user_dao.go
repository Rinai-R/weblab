package dao

import "weblab/internal/model"

type UserDAO interface {
	Create(user model.User) (model.User, error)
	GetByUsername(username string) (model.User, error)
	GetByID(id int64) (model.User, error)
	ListAll() ([]model.User, error)
	Follow(followerID, followeeID int64) error
	IsFollowing(followerID, followeeID int64) (bool, error)
	ListFollowing(userID int64) ([]int64, error)
}
