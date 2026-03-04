package handler

import (
	"weblab/internal/service"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type authReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	user, token, err := h.authSvc.Register(req.Username, req.Password)
	if err != nil {
		handleErr(c, err)
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErr(c, service.ErrInvalidArgument)
		return
	}

	user, token, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		handleErr(c, err)
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		handleErr(c, service.ErrUnauthorized)
		return
	}

	user, err := h.authSvc.Me(userID)
	if err != nil {
		handleErr(c, err)
		return
	}
	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
