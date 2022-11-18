package system

import (
	"gameoflife/utils"
	"time"
)

type SystemI interface {
	GetUserParam(userID utils.UserID) (UsersParam, error)
	GetCurrentTime() int64
	MoveCurrentTimeForvard(int64) int64
	CreateUserWithRamdomParam(user utils.UserID) error
	GetUserList() []utils.UserID
	AddUserParam(userId utils.UserID, amperes, volts, cyberlinks, kw, hydrogen int64) error
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

func (s *System) LockHydrogen(userId utils.UserID, count uint64) error {
	user, ok := s.Users[userId]
	if !ok {
		return utils.ErrSystemUserNotExist
	}
	err := user.LockHydogen(count)
	s.Users[userId] = user
	return err
}

func (s *System) AddUserParam(userId utils.UserID, amperes, volts, cyberlinks, kw, hydrogen int64) error {
	user, ok := s.Users[userId]
	if !ok {
		return utils.ErrSystemUserNotExist
	}
	user.Amperes = uint64(int64(user.Amperes) + amperes)
	user.Volts = uint64(int64(user.Volts) + volts)
	user.Cyberlinks = uint64(int64(user.Cyberlinks) + cyberlinks)
	user.Kw = uint64(int64(user.Kw) + kw)
	user.Hydrogen = uint64(int64(user.Hydrogen) + hydrogen)
	s.Users[userId] = user

	return nil
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
