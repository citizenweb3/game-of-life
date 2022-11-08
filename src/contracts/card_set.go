package contracts

import (
	"gameoflife/utils"
)

type CardWithUserInfluence struct {
	cardID         utils.CardID
	userAttributes CardParams
}

type CardSetI interface {
	GetActualSet(user utils.UserID) []utils.CardID
	AddCardToSet(executor utils.UserID, cardId utils.CardID) error
	RemoveCardFromSet(executor utils.UserID, cardId utils.CardID) error
	ChangeCardFromSet(executor utils.UserID, cardIdLast, cardIdNew utils.CardID) error

	SetUserAttribute(executor utils.UserID, numInSet uint8, value CardParams) error
	GetUserAttributes(user utils.UserID) []CardParams

	GetActualSetWithAttribute(user utils.UserID) []CardWithUserInfluence
}

type CardSet struct {
	cardsContract  CardsI
	countCardInSet uint8
	cardSet        map[utils.UserID][]utils.CardID

	personalAttributes map[utils.UserID][]CardParams
}

func NewCardSet(cardsContract CardsI, countCardInSet uint8) *CardSet {
	return &CardSet{
		cardsContract:      cardsContract,
		countCardInSet:     countCardInSet,
		cardSet:            map[utils.UserID][]utils.CardID{},
		personalAttributes: map[utils.UserID][]CardParams{},
	}
}

func (cs *CardSet) AddCardToSet(executor utils.UserID, cardId utils.CardID) error {
	err := cs.cardsContract.IsOwner(cardId, executor)
	if err != nil {
		return err
	}

	executorSet, ok := cs.cardSet[executor]
	if !ok {
		cs.cardSet[executor] = []utils.CardID{cardId}
		return nil
	}

	if ok && len(executorSet) == int(cs.countCardInSet) {
		return utils.ErrTooMuchCards
	}

	for _, cardIdInSet := range executorSet {
		if cardIdInSet == cardId {
			return utils.ErrAlreadyInSet
		}
	}

	executorSet = append(executorSet, cardId)
	cs.cardSet[executor] = executorSet

	return nil
}

func (cs *CardSet) RemoveCardFromSet(executor utils.UserID, cardId utils.CardID) error {
	executorSet, ok := cs.cardSet[executor]
	if !ok {
		return utils.ErrSetIsEmpty
	}

	for num, cardIdInSet := range executorSet {
		if cardIdInSet == cardId {
			cs.cardSet[executor] = append(executorSet[:num], executorSet[num+1:]...)
			return nil
		}
	}

	return utils.ErrCardIsNotInSet
}

func (cs *CardSet) ChangeCardFromSet(executor utils.UserID, cardIdLast, cardIdNew utils.CardID) error {
	err := cs.cardsContract.IsOwner(cardIdNew, executor)
	if err != nil {
		return err
	}

	executorSet, ok := cs.cardSet[executor]
	if !ok {
		return utils.ErrSetIsEmpty
	}

	for num, cardIdInSet := range executorSet {
		if cardIdInSet == cardIdLast {
			cs.cardSet[executor][num] = cardIdNew
			return nil
		}
	}

	return utils.ErrCardIsNotInSet
}

func (cs *CardSet) GetActualSet(user utils.UserID) []utils.CardID {
	return cs.cardSet[user]
}

///// Attributes
func (cs *CardSet) SetUserAttribute(executor utils.UserID, numInSet uint8, value CardParams) error {
	if numInSet > cs.countCardInSet {
		return utils.ErrOutOfSetRange
	}

	_, ok := cs.personalAttributes[executor]
	if !ok {
		cs.personalAttributes[executor] = make([]CardParams, int(cs.countCardInSet))
	}

	cs.personalAttributes[executor][numInSet] = value

	return nil
}

func (cs *CardSet) GetUserAttributes(user utils.UserID) []CardParams {
	return cs.personalAttributes[user]
}

func (cs *CardSet) GetActualSetWithAttribute(user utils.UserID) []CardWithUserInfluence {
	cardIDs := cs.GetActualSet(user)
	res := make([]CardWithUserInfluence, len(cardIDs))

	attributes := cs.GetUserAttributes(user)
	for num, cardID := range cardIDs {
		res[num].cardID = cardID
	}

	if len(attributes) == 0 {
		return res
	}
	for num := range cardIDs {
		res[num].userAttributes = attributes[num]
	}
	return res
}
