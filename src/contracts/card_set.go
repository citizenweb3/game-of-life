package contracts

import (
	"fmt"
	"gameoflife/utils"
)

type Influence struct {
	Hp       float64
	Level    float64
	Deffence float64
	Damage   float64
	Accuracy float64
}

func (i *Influence) Add(b Influence) {
	i.Hp += b.Hp
	i.Level += b.Level
	i.Deffence += b.Deffence
	i.Damage += b.Damage
	i.Accuracy += b.Accuracy
}

func (i *Influence) Check() bool {
	return i.Hp <= float64(100) &&
		i.Level <= float64(100) &&
		i.Deffence <= float64(100) &&
		i.Damage <= float64(100) &&
		i.Accuracy <= float64(100)
}

type CardWithUserInfluence struct {
	CardID         utils.CardID
	UserAttributes Influence
}

type CardSetI interface {
	GetActualSet(user utils.UserID) []utils.CardID
	AddCardToSet(executor utils.UserID, numInSet int, cardId utils.CardID) error
	RemoveCardFromSet(executor utils.UserID, cardId utils.CardID) error
	ChangeCardFromSet(executor utils.UserID, cardIdLast, cardIdNew utils.CardID) error

	SetUserAttribute(executor utils.UserID, numInSet uint8, value Influence) error
	GetUserAttributes(user utils.UserID) []Influence

	GetActualSetWithAttribute(user utils.UserID) []CardWithUserInfluence
}

type CardSet struct {
	cardsContract  CardsI
	countCardInSet uint8
	cardSet        map[utils.UserID][]utils.CardID

	personalAttributes map[utils.UserID][]Influence
}

func NewCardSet(cardsContract CardsI, countCardInSet uint8) *CardSet {
	return &CardSet{
		cardsContract:      cardsContract,
		countCardInSet:     countCardInSet,
		cardSet:            map[utils.UserID][]utils.CardID{},
		personalAttributes: map[utils.UserID][]Influence{},
	}
}

func (cs *CardSet) AddCardToSet(executor utils.UserID, numInSet int, cardId utils.CardID) error {
	fmt.Println("AddCardToSet to pos", numInSet, " card", cardId)
	if numInSet < 0 {
		return utils.ErrSetTooMuchCards
	}
	err := cs.cardsContract.IsOwner(cardId, executor)
	if err != nil {
		return err
	}

	executorSet, ok := cs.cardSet[executor]
	if !ok {
		cs.cardSet[executor] = make([]utils.CardID, cs.countCardInSet)
		cs.cardSet[executor][0] = cardId
		return nil
	}

	if numInSet >= int(cs.countCardInSet) {
		return utils.ErrSetTooMuchCards
	}

	for _, cardIdInSet := range executorSet {
		if cardIdInSet == cardId {
			return utils.ErrSetAlreadyInSet
		}
	}
	cs.cardSet[executor][numInSet] = cardId

	return nil
}

func (cs *CardSet) RemoveCardFromSet(executor utils.UserID, cardId utils.CardID) error {
	executorSet, ok := cs.cardSet[executor]
	if !ok {
		return utils.ErrSetIsEmpty
	}

	for num, cardIdInSet := range executorSet {
		if cardIdInSet == cardId {
			cs.cardSet[executor][num] = ""
			return nil
		}
	}

	return utils.ErrSetCardIsNotInSet
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

	numForChange := -1
	for num, cardIdInSet := range executorSet {
		if cardIdInSet == cardIdLast {
			numForChange = num
		}
		if cardIdInSet == cardIdNew {
			return utils.ErrSetAlreadyInSet
		}
	}

	if numForChange != -1 {
		cs.cardSet[executor][numForChange] = cardIdNew
		return nil
	}

	return utils.ErrSetCardIsNotInSet
}

func (cs *CardSet) GetActualSet(user utils.UserID) []utils.CardID {
	fmt.Println("GetActualSet", cs.cardSet[user])
	return cs.cardSet[user]
}

///// Attributes
func (cs *CardSet) SetUserAttribute(executor utils.UserID, numInSet uint8, value Influence) error {
	if numInSet > cs.countCardInSet {
		return utils.ErrOutOfSetRange
	}

	_, ok := cs.personalAttributes[executor]
	if !ok {
		cs.personalAttributes[executor] = make([]Influence, int(cs.countCardInSet))
	} else {
		total := value
		for _, attr := range cs.personalAttributes[executor] {
			total.Add(attr)
			if !total.Check() {
				return utils.ErrSetTooMuchCards
			}
		}
	}

	cs.personalAttributes[executor][numInSet] = value

	return nil
}

func (cs *CardSet) GetUserAttributes(user utils.UserID) []Influence {
	return cs.personalAttributes[user]
}

func (cs *CardSet) GetActualSetWithAttribute(user utils.UserID) []CardWithUserInfluence {
	cardIDs := cs.GetActualSet(user)
	res := make([]CardWithUserInfluence, cs.countCardInSet)

	if len(cardIDs) == 0 {
		return res
	}

	for num, cardID := range cardIDs {
		err := cs.cardsContract.IsOwner(cardID, user)
		if err != nil {
			continue
		}
		if cs.cardsContract.IsFreezed(cardID) {
			continue
		}
		res[num] = CardWithUserInfluence{
			CardID: cardID,
		}

	}

	for num, attributes := range cs.GetUserAttributes(user) {
		res[num].UserAttributes = attributes
	}
	return res
}
