package controllers

import (
	"gameoflife/system"
)

// battle

type SetBattleReadinessRequest struct {
	Executor string
	Ready    bool
}

type BattleRequest struct {
	Executor string
	Rival    string
}

// card set
type GetActualSetResponse []GetActualSetResponseElement
type GetActualSetResponseElement struct {
	CardID   string
	Owner    string
	Hp       uint64
	Level    uint8
	Strength uint64
	Accuracy uint64
}

type AddCardToSetRequest struct {
	Executor string
	NumInSet int
	CardID   string
}

type RemoveCardToSetRequest struct {
	Executor string
	CardID   string
}

type ChangeCardFromSetRequest struct {
	Executor   string
	CardIDLast string
	CardIDNew  string
}

type SetUserAttributeRequest struct {
	Executor string
	NumInSet uint8
	Hp       uint64
	Level    uint8
	Strength uint64
	Accuracy uint64
}

//card

type BurnRequest struct {
	Executor string
	CardID   string
}
type FreezeRequest struct {
	Executor string
	CardID1  string
	CardID2  string
}

type TransferRequest struct {
	Executor string
	CardID   string
	To       string
}

type MintNewCardRequest struct {
	UserID string
}

type MintNewCardResponse struct {
	CardId   string
	UserID   string
	Hp       uint64
	Level    uint8
	Strength uint64
	Accuracy uint64
}

type GetCardPropertiesResponce struct {
	CardId   string
	UserID   string
	Hp       uint64
	Level    uint8
	Strength uint64
	Accuracy uint64
}

// system

type MoveForwardRequest struct {
	AddUnixTime int64
}

type MoveForwardResponse struct {
	TotalMoved  int64
	CurrentTime int64
}
type AddUserParamRequest struct {
	UserID     string
	Volts      int64
	Amperes    int64
	Cyberlinks int64
	Kw         int64
}
type UserInfoApi struct {
	UserID     string
	Volts      uint64
	Amperes    uint64
	Cyberlinks uint64
	Kw         uint64
	Random     bool
}

func (uia *UserInfoApi) AddParams(up system.UsersParam) {
	uia.Volts = up.Volts
	uia.Amperes = up.Amperes
	uia.Cyberlinks = up.Cyberlinks
	uia.Kw = up.Kw
}

func (uia *UserInfoApi) GetParams() system.UsersParam {
	return system.UsersParam{
		Volts:      uia.Volts,
		Amperes:    uia.Amperes,
		Cyberlinks: uia.Cyberlinks,
		Kw:         uia.Kw,
	}
}

type GetUserListResponse []UserInfoApi

// automatisation

type CreateSystemRequest struct {
	UserCountFrom int
	UserCountTo   int

	CardCountFrom int
	CardCountTo   int
}
