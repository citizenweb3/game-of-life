package utils

import "errors"

// cards
var (
	ErrPermDenie      = errors.New("Permition denie")
	ErrCardNotExist   = errors.New("Card is not exist")
	ErrCardNotCreated = errors.New("Card has not created")
	ErrCardFreezed    = errors.New("Card is freezed")
	ErrCardNotFreezed = errors.New("Card is not freezed")
)

// cards set
var (
	ErrTooMuchCards   = errors.New("Set is already full")
	ErrSetIsEmpty     = errors.New("Set is empty")
	ErrCardIsNotInSet = errors.New("Card is not from set")
	ErrAlreadyInSet   = errors.New("Card already in set")
	ErrOutOfSetRange  = errors.New("Out of set range")
)

// battle
var (
	ErrUserIsNotReady = errors.New("Rival is not ready to battle")
)

// system
var (
	ErrUserAlreadyExist = errors.New("User already exist")
	ErrUserNotExist     = errors.New("User not exist")
)
