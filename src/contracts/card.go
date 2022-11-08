package contracts

import (
	"math/rand"
	"time"

	"gameoflife/utils"
)

// AddAvatar(cardId utils.CardID, url string, executor utils.UserID) error
type CardsI interface {
	Transfer(cardId utils.CardID, executor, to utils.UserID) error
	Burn(cardId utils.CardID, executor utils.UserID) error
	Freeze(cardId utils.CardID, executor utils.UserID) error
	UnFreeze(cardId utils.CardID, executor utils.UserID) error
	MintNewCard(executor utils.UserID) (utils.CardID, error)

	GetCardProperties(cardId utils.CardID) (CardParams, error)
	GetCardOwner(cardId utils.CardID) (utils.UserID, error)

	GetOwnersCards(executor utils.UserID) []Card
	IsOwner(cardId utils.CardID, executor utils.UserID) error

	IsFreezed(cardId utils.CardID) bool
}

type CardParams struct {
	Hp       uint64
	Level    uint8
	Strength uint64
	Accuracy uint64
}

type Card struct {
	Id     utils.CardID
	Params CardParams
}

type Cards struct {
	ownerCards  map[utils.UserID][]utils.CardID
	cards       map[utils.CardID]Card
	cardOwner   map[utils.CardID]utils.UserID
	freezedCard map[utils.CardID]uint64
	freezeTime  uint64
}

func NewCards(freezeTime uint64) *Cards {
	return &Cards{
		ownerCards:  map[utils.UserID][]utils.CardID{},
		cards:       map[utils.CardID]Card{},
		cardOwner:   map[utils.CardID]utils.UserID{},
		freezedCard: map[utils.CardID]uint64{},
		freezeTime:  freezeTime,
	}
}
func (c *Cards) GetCardOwner(cardId utils.CardID) (utils.UserID, error) {
	card, exist := c.cardOwner[cardId]
	if !exist {
		return "", utils.ErrCardNotExist
	}
	return card, nil
}

func (c *Cards) GetOwnersCards(executor utils.UserID) []Card {
	cardIds, ok := c.ownerCards[executor]
	if !ok {
		return []Card{}
	}
	resp := make([]Card, 0, len(cardIds))

	for _, id := range cardIds {
		card := c.cards[id]
		resp = append(resp, card)
	}
	return resp
}

func (c *Cards) GetCardProperties(cardId utils.CardID) (CardParams, error) {
	card, exist := c.cards[cardId]
	if !exist {
		return CardParams{}, utils.ErrCardNotExist
	}
	return card.Params, nil
}

func (c *Cards) MintNewCard(executor utils.UserID) (utils.CardID, error) {
	var cp CardParams
	var exist bool
	var hash utils.CardID
	for i := 0; i < 5; i++ {
		cp = generarateRandomCardParams()
		hash = utils.CardID(utils.Hash(cp))
		_, exist = c.cards[hash]
		if !exist {
			break
		}
	}
	if exist {
		return "", utils.ErrCardNotCreated
	}

	c.cardOwner[hash] = executor
	c.cards[hash] = Card{
		Id:     hash,
		Params: cp,
	}
	if executorCards, ok := c.ownerCards[executor]; ok {
		executorCards = append(executorCards, hash)
		c.ownerCards[executor] = executorCards
	} else {
		c.ownerCards[executor] = []utils.CardID{hash}
	}
	return hash, nil
}

func generarateRandomCardParams() CardParams {
	return CardParams{
		Hp:       rand.Uint64() % 100,
		Level:    uint8(rand.Uint32() % 5),
		Strength: rand.Uint64() % 10,
		Accuracy: rand.Uint64() % 100,
	}
}

func (c *Cards) IsOwner(cardId utils.CardID, executor utils.UserID) error {
	owner, ok := c.cardOwner[cardId]
	if !ok {
		return utils.ErrCardNotExist
	}

	if owner != executor {
		return utils.ErrPermDenie
	}
	return nil
}

func (c *Cards) Transfer(cardId utils.CardID, executor, to utils.UserID) error {
	if err := c.IsOwner(cardId, executor); err != nil {
		return err
	}

	c.cardOwner[cardId] = to
	cards := c.ownerCards[executor]

	for num, val := range cards {
		if val == cardId {
			c.ownerCards[executor] = append(cards[:num], cards[num+1:]...)
			break
		}
	}

	cards2 := c.ownerCards[to]
	cards2 = append(cards2, cardId)
	c.ownerCards[to] = cards2

	return nil
}

func (c *Cards) Burn(cardId utils.CardID, executor utils.UserID) error {
	if err := c.IsOwner(cardId, executor); err != nil {
		return err
	}

	cards := c.ownerCards[executor]
	delete(c.cardOwner, cardId)

	for num, val := range cards {
		if val == cardId {
			c.ownerCards[executor] = append(cards[:num], cards[num+1:]...)
			break
		}
	}

	return nil
}

func (c *Cards) IsFreezed(cardId utils.CardID) bool {
	if _, ok := c.freezedCard[cardId]; ok {
		return true
	}
	return false
}

func (c *Cards) Freeze(cardId utils.CardID, executor utils.UserID) error {
	if err := c.IsOwner(cardId, executor); err != nil {
		return err
	}

	if c.IsFreezed(cardId) {
		return utils.ErrCardFreezed
	}

	c.freezedCard[cardId] = c.freezeTime + uint64(time.Now().Unix())

	return nil
}

func (c *Cards) UnFreeze(cardId utils.CardID, executor utils.UserID) error {
	if err := c.IsOwner(cardId, executor); err != nil {
		return err
	}

	if !c.IsFreezed(cardId) {
		return utils.ErrCardNotFreezed
	}

	finishFreezeTime := c.freezedCard[cardId]

	if finishFreezeTime > uint64(time.Now().Unix()) {
		return utils.ErrCardFreezed
	}

	delete(c.freezedCard, cardId)
	// todo add rewards to freeze

	return nil
}
