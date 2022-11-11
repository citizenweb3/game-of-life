package system

import (
	"gameoflife/utils"
	"time"
)

type SystemI interface {
	GetUserParam(userID utils.UserID) (UsersParam, error)
	GetCurrentTime() int64
	MoveCurrentTimeForvard(int64) int64
}

type System struct {
	Users           map[utils.UserID]UsersParam
	moveForvardTime int64
}

func NewSystem() *System {
	return &System{
		Users:           map[utils.UserID]UsersParam{},
		moveForvardTime: 0,
	}
}

func (s *System) CreateUserWithRamdomParam(user utils.UserID) error {
	return s.CreateUserWithParam(user, GenerateRandomUserParam())
}

func (s *System) CreateUserWithParam(user utils.UserID, up UsersParam) error {
	if _, ok := s.Users[user]; ok {
		return utils.ErrSystemUserAlreadyExist
	}
	s.Users[user] = up
	return nil
}

func (s *System) GetUserParam(userID utils.UserID) (UsersParam, error) {
	user, ok := s.Users[userID]
	if !ok {
		return user, utils.ErrSystemUserNotExist
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

func (s *System) GetCurrentTime() int64 {
	return time.Now().Unix() + s.moveForvardTime
}

func (s *System) MoveCurrentTimeForvard(moveTo int64) int64 {
	s.moveForvardTime += moveTo
	return s.moveForvardTime
}
