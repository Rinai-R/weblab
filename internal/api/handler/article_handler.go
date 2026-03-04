package handler

import (
	"weblab/internal/service"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleSvc *service.ArticleService
}

func NewArticleHandler(articleSvc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleSvc: articleSvc}
}

type publishArticleReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *ArticleHandler) Publish(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	var req publishArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	article, err := h.articleSvc.Publish(userID, req.Title, req.Content)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, article)
}

func (h *ArticleHandler) GetDetail(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	id, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	detail, err := h.articleSvc.GetDetail(userID, id)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, detail)
}

func (h *ArticleHandler) Recommend(c *gin.Context) {
	limit, err := queryInt(c, "limit", 20)
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	feed, err := h.articleSvc.Recommend(limit)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, feed)
}

func (h *ArticleHandler) Feed(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	feed, err := h.articleSvc.Feed(userID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, feed)
}
