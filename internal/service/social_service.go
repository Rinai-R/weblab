package service

import (
	"strings"
	"weblab/internal/dao"
	"weblab/internal/model"
)

type SocialService struct {
	userDAO    dao.UserDAO
	messageDAO dao.MessageDAO
}

func NewSocialService(userDAO dao.UserDAO, messageDAO dao.MessageDAO) *SocialService {
	return &SocialService{userDAO: userDAO, messageDAO: messageDAO}
}

func (s *SocialService) Follow(followerID, followeeID int64) error {
	if followerID <= 0 || followeeID <= 0 || followerID == followeeID {
		return ErrInvalidArgument
	}
	if _, err := s.userDAO.GetByID(followerID); err != nil {
		return err
	}
	if _, err := s.userDAO.GetByID(followeeID); err != nil {
		return err
	}
	return s.userDAO.Follow(followerID, followeeID)
}

func (s *SocialService) SendMessage(fromUserID, toUserID int64, content string) (model.Message, error) {
	content = strings.TrimSpace(content)
	if fromUserID <= 0 || toUserID <= 0 || content == "" {
		return model.Message{}, ErrInvalidArgument
	}

	if err := s.ensureMutualFollow(fromUserID, toUserID); err != nil {
		return model.Message{}, err
	}

	return s.messageDAO.Create(model.Message{FromUserID: fromUserID, ToUserID: toUserID, Content: content})
}

func (s *SocialService) Conversation(userID, peerID int64) ([]model.Message, error) {
	if userID <= 0 || peerID <= 0 {
		return nil, ErrInvalidArgument
	}

	if err := s.ensureMutualFollow(userID, peerID); err != nil {
		return nil, err
	}

	return s.messageDAO.ListConversation(userID, peerID)
}

func (s *SocialService) DiscoverUsers(userID int64) ([]model.UserRelation, error) {
	if userID <= 0 {
		return nil, ErrInvalidArgument
	}
	if _, err := s.userDAO.GetByID(userID); err != nil {
		return nil, err
	}

	users, err := s.userDAO.ListAll()
	if err != nil {
		return nil, err
	}

	relations := make([]model.UserRelation, 0, len(users))
	for _, u := range users {
		if u.ID == userID {
			continue
		}

		isFollowing, err := s.userDAO.IsFollowing(userID, u.ID)
		if err != nil {
			return nil, err
		}
		isFollowedByPeer, err := s.userDAO.IsFollowing(u.ID, userID)
		if err != nil {
			return nil, err
		}

		relations = append(relations, model.UserRelation{
			ID:             u.ID,
			Username:       u.Username,
			IsFollowing:    isFollowing,
			IsMutualFollow: isFollowing && isFollowedByPeer,
		})
	}
	return relations, nil
}

func (s *SocialService) MutualFollows(userID int64) ([]model.UserBrief, error) {
	relations, err := s.DiscoverUsers(userID)
	if err != nil {
		return nil, err
	}

	peers := make([]model.UserBrief, 0)
	for _, relation := range relations {
		if !relation.IsMutualFollow {
			continue
		}
		peers = append(peers, model.UserBrief{
			ID:       relation.ID,
			Username: relation.Username,
		})
	}
	return peers, nil
}

func (s *SocialService) ensureMutualFollow(a, b int64) error {
	if _, err := s.userDAO.GetByID(a); err != nil {
		return err
	}
	if _, err := s.userDAO.GetByID(b); err != nil {
		return err
	}

	ab, err := s.userDAO.IsFollowing(a, b)
	if err != nil {
		return err
	}
	ba, err := s.userDAO.IsFollowing(b, a)
	if err != nil {
		return err
	}
	if !ab || !ba {
		return ErrForbidden
	}
	return nil
}
