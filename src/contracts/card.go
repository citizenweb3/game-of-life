package contracts

import (
	"net/url"

	"gameoflife/system"
	"gameoflife/utils"
)

var MIXED_ADD = uint64(5)

const (
	HP_RANGE = uint64(100)
	HP_MIN   = uint64(100)

	ACCURACY_RANGE = uint64(20)
	ACCURACY_MIN   = uint64(10)

	LEVEL_RANGE = uint64(5)
	LEVEL_MIN   = uint64(0)

	DAMAGE_RANGE = uint64(30)
	DAMAGE_MIN   = uint64(20)

	DEFFENCE_RANGE = uint64(20)
	DEFFENCE_MIN   = uint64(0)
)

// AddAvatar(cardId utils.CardID, url string, executor utils.UserID) error
type CardsI interface {
	Transfer(cardId utils.CardID, executor, to utils.UserID) error
	Burn(cardId utils.CardID, executor utils.UserID) error
	Freeze(cardId1, cardId2 utils.CardID, executor utils.UserID) error
	GetFreezeTime(cardId utils.CardID) int64
	UnFreeze(cardId utils.CardID, executor utils.UserID) error
	MintNewCard(executor utils.UserID) (utils.CardID, error)

	GetCardProperties(cardId utils.CardID) (CardParams, error)
	GetCardOwner(cardId utils.CardID) (utils.UserID, error)

	GetOwnersCards(executor utils.UserID) []Card
	IsOwner(cardId utils.CardID, executor utils.UserID) error

	IsFreezed(cardId utils.CardID) bool
}

type CardParams struct {
	Hp    uint64 // Kw
	Level uint64 // CyberLinks

	Deffence uint64 // locked hydrogen

	Damage   uint64 // Ampers
	Accuracy uint64 // Volts

}

func (cp *CardParams) ByUserParam(up system.UsersParam) {
	cp.Hp = up.GetKw()
	cp.Level = up.GetCountCyberlinks()
	cp.Deffence = up.GetLockedHydrogen()
	cp.Damage = up.GetAmperes()
	cp.Accuracy = up.GetVolts()
}

type Card struct {
	Id     utils.CardID
	Params CardParams
	Avatar url.URL
}

type Cards struct {
	system system.SystemI

	ownerCards   map[utils.UserID][]utils.CardID
	cards        map[utils.CardID]Card
	cardOwner    map[utils.CardID]utils.UserID
	freezedPair  map[utils.CardID]utils.CardID
	freezedUntil map[string]int64
	freezeTime   int64
}

func NewCards(system system.SystemI, freezeTime uint32) *Cards {
	return &Cards{
		system:       system,
		ownerCards:   map[utils.UserID][]utils.CardID{},
		cards:        map[utils.CardID]Card{},
		cardOwner:    map[utils.CardID]utils.UserID{},
		freezedPair:  map[utils.CardID]utils.CardID{},
		freezedUntil: map[string]int64{},
		freezeTime:   int64(freezeTime),
	}
}

func (c *Cards) GetFreezeTime(cardId utils.CardID) int64 {
	timeFreeze, ok := c.freezedUntil[cardId.ToString()]
	if ok {
		return timeFreeze
	}
	freezedPair, ok := c.freezedPair[cardId]
	if !ok {
		return 0
	}
	timeFreeze, ok = c.freezedUntil[freezedPair.ToString()]
	if !ok {
		return 0
	}
	return timeFreeze
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
	var err error
	var hash utils.CardID
	for i := 0; i < 5; i++ {
		cp, err = c.GenerarateRandomCardParams(executor)
		if err != nil {
			continue
		}
		hash = utils.CardID(utils.Hash(cp))
		_, exist = c.cards[hash]
		if !exist {
			break
		}
	}
	if exist {
		return "", utils.ErrCardNotCreated
	}
	card := Card{
		Id:     hash,
		Params: cp,
	}
	c.addCard(card, executor)

	return hash, nil
}

func (c *Cards) addCard(card Card, owner utils.UserID) {
	c.cardOwner[card.Id] = owner
	c.cards[card.Id] = card
	if executorCards, ok := c.ownerCards[owner]; ok {
		executorCards = append(executorCards, card.Id)
		c.ownerCards[owner] = executorCards
	} else {
		c.ownerCards[owner] = []utils.CardID{card.Id}
	}
}
func (c *Cards) GenerarateRandomCardParams(userID utils.UserID) (CardParams, error) {
	userParam, err := c.system.GetUserParam(userID)
	var cardParam CardParams
	if err != nil {
		return cardParam, err
	}
	cardParam = generarateRandomCardParamsWithLevel(1 + utils.GetRandomNumberUint64(userParam.GetCountCyberlinks()/100))

	return cardParam, nil
}

func generarateRandomCardParamsWithLevel(level uint64) CardParams {
	return CardParams{
		Hp:    (utils.GetRandomNumberUint64(HP_RANGE) + HP_MIN) * level,
		Level: level,

		Deffence: (utils.GetRandomNumberUint64(DEFFENCE_RANGE) + DEFFENCE_MIN) * level,

		Damage:   (utils.GetRandomNumberUint64(DAMAGE_RANGE) + DAMAGE_MIN) * level,
		Accuracy: (utils.GetRandomNumberUint64(ACCURACY_RANGE) + ACCURACY_MIN) * level,
	}

}
func generarateRandomCardParams() CardParams {
	return CardParams{
		Hp:       utils.GetRandomNumberUint64(100),
		Level:    utils.GetRandomNumberUint64(5),
		Accuracy: utils.GetRandomNumberUint64(100),
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

	c.unsavedBurn(cardId, executor)

	return nil
}
func (c *Cards) unsavedBurn(cardId utils.CardID, owner utils.UserID) {
	cards := c.ownerCards[owner]
	delete(c.cardOwner, cardId)
	delete(c.cards, cardId)

	for num, val := range cards {
		if val == cardId {
			c.ownerCards[owner] = append(cards[:num], cards[num+1:]...)
			break
		}
	}
}

func (c *Cards) IsFreezed(cardId utils.CardID) bool {
	if _, ok := c.freezedPair[cardId]; ok {
		return true
	}
	return false
}

func (c *Cards) Freeze(cardId1, cardId2 utils.CardID, executor utils.UserID) error {
	if cardId1 == cardId2 {
		return utils.ErrCardsEqual
	}

	if err := c.IsOwner(cardId1, executor); err != nil {
		return err
	}
	if err := c.IsOwner(cardId2, executor); err != nil {
		return err
	}

	if c.IsFreezed(cardId1) {
		return utils.ErrCardFreezed
	}
	if c.IsFreezed(cardId2) {
		return utils.ErrCardFreezed
	}

	c.freezedPair[cardId1] = cardId2
	c.freezedPair[cardId2] = cardId1

	c.freezedUntil[cardId1.ToString()] = c.freezeTime + c.system.GetCurrentTime()

	return nil
}

func (c *Cards) UnFreeze(cardId utils.CardID, executor utils.UserID) error {
	if err := c.IsOwner(cardId, executor); err != nil {
		return err
	}

	if !c.IsFreezed(cardId) {
		return utils.ErrCardNotFreezed
	}

	finishFreezeTime, ok := c.freezedUntil[cardId.ToString()]
	if ok {
		if finishFreezeTime > c.system.GetCurrentTime() {
			return utils.ErrCardFreezed
		}
		delete(c.freezedUntil, cardId.ToString())

	}
	freezePairCard := c.freezedPair[cardId]
	if !ok {
		finishFreezeTime = c.freezedUntil[freezePairCard.ToString()]
		if finishFreezeTime > c.system.GetCurrentTime() {
			return utils.ErrCardFreezed
		}
		delete(c.freezedUntil, freezePairCard.ToString())
	}

	delete(c.freezedPair, cardId)
	delete(c.freezedPair, freezePairCard)

	mixedCard := c.mixedCards(cardId, freezePairCard)

	c.unsavedBurn(cardId, executor)
	c.unsavedBurn(freezePairCard, executor)
	c.addCard(mixedCard, executor)

	return nil
}

// mixedCards - mixed 2 cards. Attantion: card must exist!!!
func (c *Cards) mixedCards(cardId1, cardId2 utils.CardID) (res Card) {
	cardParam1, _ := c.GetCardProperties(cardId1)
	cardParam2, _ := c.GetCardProperties(cardId2)

	res.Params.Hp = getRandBeetweenUint64(cardParam1.Hp, cardParam2.Hp)
	res.Params.Accuracy = getRandBeetweenUint64(cardParam1.Accuracy, cardParam2.Accuracy)
	if cardParam1.Level > cardParam2.Level {
		res.Params.Level = cardParam1.Level
	} else {
		res.Params.Level = cardParam2.Level
	}
	res.Params.Damage = getRandBeetweenUint64(cardParam1.Damage, cardParam2.Damage)
	res.Params.Deffence = getRandBeetweenUint64(cardParam1.Deffence, cardParam2.Deffence)
	res.Id = utils.CardID(utils.Hash(res.Params))
	return
}

func getRandBeetweenUint64(p1, p2 uint64) uint64 {
	if p1 > p2 {
		p2, p1 = p1, p2
	}
	diff := p2 - p1
	return p1 + utils.GetRandomNumberUint64(diff+MIXED_ADD)
}
