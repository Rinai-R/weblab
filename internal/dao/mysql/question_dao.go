package mysql

import (
	"time"
	"weblab/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuestionDAO struct {
	db *gorm.DB
}

func NewQuestionDAO(db *gorm.DB) *QuestionDAO {
	return &QuestionDAO{db: db}
}

func (d *QuestionDAO) CreateQuestion(question model.Question) (model.Question, error) {
	record := QuestionRecord{
		AuthorID:    question.AuthorID,
		Title:       question.Title,
		Description: question.Description,
	}
	if err := d.db.Create(&record).Error; err != nil {
		return model.Question{}, normalizeErr(err)
	}

	return model.Question{
		ID:          record.ID,
		AuthorID:    record.AuthorID,
		Title:       record.Title,
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
	}, nil
}

func (d *QuestionDAO) GetQuestionByID(id int64) (model.Question, error) {
	var record QuestionRecord
	if err := d.db.Where("id = ?", id).First(&record).Error; err != nil {
		return model.Question{}, normalizeErr(err)
	}

	return model.Question{
		ID:          record.ID,
		AuthorID:    record.AuthorID,
		Title:       record.Title,
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
	}, nil
}

func (d *QuestionDAO) ListQuestions(limit int, cursor int64) ([]model.Question, error) {
	if limit <= 0 {
		limit = 20
	}

	var records []QuestionRecord
	query := d.db
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}

	err := query.Order("created_at DESC").Order("id DESC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	questions := make([]model.Question, 0, len(records))
	for _, record := range records {
		questions = append(questions, model.Question{
			ID:          record.ID,
			AuthorID:    record.AuthorID,
			Title:       record.Title,
			Description: record.Description,
			CreatedAt:   record.CreatedAt,
		})
	}
	return questions, nil
}

func (d *QuestionDAO) ListQuestionsByAuthorIDs(authorIDs []int64, limit int, cursor int64) ([]model.Question, error) {
	if len(authorIDs) == 0 {
		return []model.Question{}, nil
	}
	if limit <= 0 {
		limit = 20
	}

	var records []QuestionRecord
	query := d.db.Where("author_id IN ?", authorIDs)
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}

	err := query.Order("created_at DESC").Order("id DESC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	questions := make([]model.Question, 0, len(records))
	for _, record := range records {
		questions = append(questions, model.Question{
			ID:          record.ID,
			AuthorID:    record.AuthorID,
			Title:       record.Title,
			Description: record.Description,
			CreatedAt:   record.CreatedAt,
		})
	}
	return questions, nil
}

func (d *QuestionDAO) CreateAnswer(answer model.Answer) (model.Answer, error) {
	record := AnswerRecord{
		QuestionID: answer.QuestionID,
		AuthorID:   answer.AuthorID,
		Content:    answer.Content,
	}
	if err := d.db.Create(&record).Error; err != nil {
		return model.Answer{}, normalizeErr(err)
	}

	return model.Answer{
		ID:         record.ID,
		QuestionID: record.QuestionID,
		AuthorID:   record.AuthorID,
		Content:    record.Content,
		CreatedAt:  record.CreatedAt,
	}, nil
}

func (d *QuestionDAO) GetAnswerByID(id int64) (model.Answer, error) {
	var record AnswerRecord
	if err := d.db.Where("id = ?", id).First(&record).Error; err != nil {
		return model.Answer{}, normalizeErr(err)
	}

	return model.Answer{
		ID:         record.ID,
		QuestionID: record.QuestionID,
		AuthorID:   record.AuthorID,
		Content:    record.Content,
		CreatedAt:  record.CreatedAt,
	}, nil
}

func (d *QuestionDAO) ListAnswers(questionID int64, limit int, cursor int64) ([]model.Answer, error) {
	if limit <= 0 {
		limit = 20
	}

	var records []AnswerRecord
	query := d.db.Where("question_id = ?", questionID)
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}

	err := query.Order("created_at DESC").Order("id DESC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, normalizeErr(err)
	}

	answers := make([]model.Answer, 0, len(records))
	for _, record := range records {
		answers = append(answers, model.Answer{
			ID:         record.ID,
			QuestionID: record.QuestionID,
			AuthorID:   record.AuthorID,
			Content:    record.Content,
			CreatedAt:  record.CreatedAt,
		})
	}
	return answers, nil
}

func (d *QuestionDAO) CountAnswers(questionID int64) (int, error) {
	var count int64
	err := d.db.Model(&AnswerRecord{}).Where("question_id = ?", questionID).Count(&count).Error
	if err != nil {
		return 0, normalizeErr(err)
	}
	return int(count), nil
}

func (d *QuestionDAO) VoteAnswer(userID, answerID int64) error {
	now := time.Now()
	record := AnswerVoteRecord{
		AnswerID:  answerID,
		UserID:    userID,
		CreatedAt: now,
	}

	err := d.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "answer_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"created_at": now,
		}),
	}).Create(&record).Error
	return normalizeErr(err)
}

func (d *QuestionDAO) UnvoteAnswer(userID, answerID int64) error {
	err := d.db.Where("answer_id = ? AND user_id = ?", answerID, userID).Delete(&AnswerVoteRecord{}).Error
	return normalizeErr(err)
}

func (d *QuestionDAO) IsAnswerVoted(userID, answerID int64) (bool, error) {
	var count int64
	err := d.db.Model(&AnswerVoteRecord{}).Where("answer_id = ? AND user_id = ?", answerID, userID).Count(&count).Error
	if err != nil {
		return false, normalizeErr(err)
	}
	return count > 0, nil
}

func (d *QuestionDAO) CountAnswerVotes(answerID int64) (int, error) {
	var count int64
	err := d.db.Model(&AnswerVoteRecord{}).Where("answer_id = ?", answerID).Count(&count).Error
	if err != nil {
		return 0, normalizeErr(err)
	}
	return int(count), nil
}
