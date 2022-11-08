package system

import (
	"gameoflife/utils"
)

type SystemI interface {
	GetUserParam(userID utils.UserID) (UsersParam, error)
}

type System struct {
	Users map[utils.UserID]UsersParam
}

func NewSystem() *System {
	return &System{
		Users: map[utils.UserID]UsersParam{},
	}
}

func (s *System) CreateUserWithRamdomParam(user utils.UserID) error {
	return s.CreateUserWithParam(user, GenerateRandomUserParam())
}

func (s *System) CreateUserWithParam(user utils.UserID, up UsersParam) error {
	if _, ok := s.Users[user]; ok {
		return utils.ErrUserAlreadyExist
	}
	s.Users[user] = up
	return nil
}

func (s *System) GetUserParam(userID utils.UserID) (UsersParam, error) {
	user, ok := s.Users[userID]
	if !ok {
		return user, utils.ErrUserNotExist
	}

	return user, nil
}

func (s *System) GetUserList() []utils.UserID {
	res := make([]utils.UserID, 0, len(s.Users))
	for userId := range s.Users {
		res = append(res, userId)
	}
	return res
}
