package utils

import "errors"

// cards
var (
	ErrPermDenie      = errors.New("Permition denie")
	ErrCardNotExist   = errors.New("Card is not exist")
	ErrCardNotCreated = errors.New("Card has not created")
	ErrCardFreezed    = errors.New("Card is freezed")
	ErrCardNotFreezed = errors.New("Card is not freezed")
	ErrCardsEqual     = errors.New("Cards are equal")
)

// cards set
var (
	ErrSetTooMuchCards   = errors.New("Set is already full")
	ErrSetIsEmpty        = errors.New("Set is empty")
	ErrSetCardIsNotInSet = errors.New("Card is not from set")
	ErrSetAlreadyInSet   = errors.New("Card already in set")
	ErrOutOfSetRange     = errors.New("Out of set range")
)

// battle
var (
	ErrBattleUserIsNotReady         = errors.New("Rival is not ready to battle")
	ErrBattleUserHasNoOneCardsInSet = errors.New("User has no one cards in set")
)

// system
var (
	ErrSystemUserAlreadyExist = errors.New("User already exist")
	ErrSystemUserNotExist     = errors.New("User not exist")
)
