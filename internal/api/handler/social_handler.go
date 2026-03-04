package handler

import (
	"weblab/internal/service"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialSvc *service.SocialService
}

func NewSocialHandler(socialSvc *service.SocialService) *SocialHandler {
	return &SocialHandler{socialSvc: socialSvc}
}

type sendMessageReq struct {
	ToUserID int64  `json:"to_user_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func (h *SocialHandler) Follow(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	followeeID, err := pathInt64(c, "id")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	if err := h.socialSvc.Follow(userID, followeeID); err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, gin.H{"followee_id": followeeID})
}

func (h *SocialHandler) SendMessage(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	var req sendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	msg, err := h.socialSvc.SendMessage(userID, req.ToUserID, req.Content)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, msg)
}

func (h *SocialHandler) Conversation(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	peerID, err := pathInt64(c, "userID")
	if err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	messages, err := h.socialSvc.Conversation(userID, peerID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, messages)
}

func (h *SocialHandler) Discover(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	users, err := h.socialSvc.DiscoverUsers(userID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, users)
}

func (h *SocialHandler) Mutuals(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	users, err := h.socialSvc.MutualFollows(userID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, users)
}
