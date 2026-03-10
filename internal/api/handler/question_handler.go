package handler

import (
	"weblab/internal/model"
	"weblab/internal/service"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionSvc *service.QuestionService
}

func NewQuestionHandler(questionSvc *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionSvc: questionSvc}
}

type askQuestionReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type answerQuestionReq struct {
	Content string `json:"content" binding:"required"`
}

func (h *QuestionHandler) Ask(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	var req askQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	question, err := h.questionSvc.Ask(userID, req.Title, req.Description)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, question)
}

func (h *QuestionHandler) Recommend(c *gin.Context) {
	limit, err := queryInt(c, "limit", 20)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}
	cursor, err := queryInt64(c, "cursor", 0)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	items, err := h.questionSvc.Recommend(limit, cursor)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, pullResp(items, questionCursor(items)))
}

func (h *QuestionHandler) FollowFeed(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	limit, err := queryInt(c, "limit", 20)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}
	cursor, err := queryInt64(c, "cursor", 0)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	items, err := h.questionSvc.FollowFeed(userID, limit, cursor)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, pullResp(items, questionCursor(items)))
}

func (h *QuestionHandler) Detail(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	questionID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	detail, err := h.questionSvc.GetDetail(userID, questionID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, detail)
}

func (h *QuestionHandler) Answer(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	questionID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	var req answerQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	answer, err := h.questionSvc.Answer(userID, questionID, req.Content)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, answer)
}

func (h *QuestionHandler) ListAnswers(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	questionID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}
	limit, err := queryInt(c, "limit", 20)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}
	cursor, err := queryInt64(c, "cursor", 0)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	answers, err := h.questionSvc.ListAnswers(userID, questionID, limit, cursor)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, pullResp(answers, answerCursor(answers)))
}

func (h *QuestionHandler) VoteAnswer(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	answerID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	if err := h.questionSvc.VoteAnswer(userID, answerID); err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, gin.H{"voted": true})
}

func (h *QuestionHandler) UnvoteAnswer(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	answerID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	if err := h.questionSvc.UnvoteAnswer(userID, answerID); err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, gin.H{"voted": false})
}

func pullResp(items interface{}, nextCursor int64) gin.H {
	return gin.H{
		"items":       items,
		"next_cursor": nextCursor,
	}
}

func questionCursor(items []model.QuestionCard) int64 {
	if len(items) == 0 {
		return 0
	}
	return items[len(items)-1].ID
}

func answerCursor(items []model.AnswerView) int64 {
	if len(items) == 0 {
		return 0
	}
	return items[len(items)-1].ID
}
