package handler

import (
	"weblab/internal/service"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	interactionSvc *service.InteractionService
}

func NewInteractionHandler(interactionSvc *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{interactionSvc: interactionSvc}
}

type addCommentReq struct {
	Content string `json:"content" binding:"required"`
}

func (h *InteractionHandler) LikeArticle(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}
	articleID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	if err := h.interactionSvc.LikeArticle(userID, articleID); err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, gin.H{"liked": true})
}

func (h *InteractionHandler) UnlikeArticle(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}
	articleID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	if err := h.interactionSvc.UnlikeArticle(userID, articleID); err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, gin.H{"liked": false})
}

func (h *InteractionHandler) AddComment(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}
	articleID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	var req addCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	comment, err := h.interactionSvc.AddComment(userID, articleID, req.Content)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, comment)
}

func (h *InteractionHandler) ListComments(c *gin.Context) {
	articleID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	comments, err := h.interactionSvc.ListComments(articleID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, comments)
}
