package service

import (
	"strings"
	"weblab/internal/dao"
	"weblab/internal/model"
	"weblab/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userDAO dao.UserDAO
	jwtMgr  *utils.JWTManager
}

func NewAuthService(userDAO dao.UserDAO, jwtMgr *utils.JWTManager) *AuthService {
	return &AuthService{userDAO: userDAO, jwtMgr: jwtMgr}
}

func (s *AuthService) Register(username, password string) (model.User, string, error) {
	username = strings.TrimSpace(username)
	if username == "" || len(password) < 6 {
		return model.User{}, "", ErrInvalidArgument
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, "", err
	}

	created, err := s.userDAO.Create(model.User{
		Username:     username,
		PasswordHash: string(hash),
	})
	if err != nil {
		return model.User{}, "", err
	}

	token, err := s.jwtMgr.Generate(created.ID)
	if err != nil {
		return model.User{}, "", err
	}
	return created, token, nil
}

func (s *AuthService) Login(username, password string) (model.User, string, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return model.User{}, "", ErrInvalidArgument
	}

	user, err := s.userDAO.GetByUsername(username)
	if err != nil {
		return model.User{}, "", ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return model.User{}, "", ErrUnauthorized
	}

	token, err := s.jwtMgr.Generate(user.ID)
	if err != nil {
		return model.User{}, "", err
	}
	return user, token, nil
}

func (s *AuthService) Me(userID int64) (model.User, error) {
	if userID <= 0 {
		return model.User{}, ErrUnauthorized
	}
	return s.userDAO.GetByID(userID)
}
