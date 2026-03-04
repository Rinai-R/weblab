package mysql

import (
	"time"
	"weblab/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) Create(user model.User) (model.User, error) {
	record := UserRecord{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}
	if err := d.db.Create(&record).Error; err != nil {
		return model.User{}, normalizeErr(err)
	}

	return model.User{
		ID:           record.ID,
		Username:     record.Username,
		PasswordHash: record.PasswordHash,
		CreatedAt:    record.CreatedAt,
	}, nil
}

func (d *UserDAO) GetByUsername(username string) (model.User, error) {
	var record UserRecord
	if err := d.db.Where("username = ?", username).First(&record).Error; err != nil {
		return model.User{}, normalizeErr(err)
	}

	return model.User{
		ID:           record.ID,
		Username:     record.Username,
		PasswordHash: record.PasswordHash,
		CreatedAt:    record.CreatedAt,
	}, nil
}

func (d *UserDAO) GetByID(id int64) (model.User, error) {
	var record UserRecord
	if err := d.db.Where("id = ?", id).First(&record).Error; err != nil {
		return model.User{}, normalizeErr(err)
	}

	return model.User{
		ID:           record.ID,
		Username:     record.Username,
		PasswordHash: record.PasswordHash,
		CreatedAt:    record.CreatedAt,
	}, nil
}

func (d *UserDAO) ListAll() ([]model.User, error) {
	var records []UserRecord
	if err := d.db.Order("id ASC").Find(&records).Error; err != nil {
		return nil, normalizeErr(err)
	}

	users := make([]model.User, 0, len(records))
	for _, record := range records {
		users = append(users, model.User{
			ID:           record.ID,
			Username:     record.Username,
			PasswordHash: record.PasswordHash,
			CreatedAt:    record.CreatedAt,
		})
	}
	return users, nil
}

func (d *UserDAO) Follow(followerID, followeeID int64) error {
	now := time.Now()
	record := FollowRecord{
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt:  now,
	}

	err := d.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "follower_id"}, {Name: "followee_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"created_at": now,
		}),
	}).Create(&record).Error
	return normalizeErr(err)
}

func (d *UserDAO) IsFollowing(followerID, followeeID int64) (bool, error) {
	var count int64
	err := d.db.Model(&FollowRecord{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	if err != nil {
		return false, normalizeErr(err)
	}
	return count > 0, nil
}

func (d *UserDAO) ListFollowing(userID int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := d.db.Model(&FollowRecord{}).
		Where("follower_id = ?", userID).
		Order("followee_id ASC").
		Pluck("followee_id", &ids).Error
	if err != nil {
		return nil, normalizeErr(err)
	}
	return ids, nil
}
