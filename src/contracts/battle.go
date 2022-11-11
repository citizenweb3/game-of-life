package contracts

import (
	"gameoflife/system"
	"gameoflife/utils"
)

type BattleI interface {
	ReadyToBattle(executor utils.UserID)
	NotReadyToBattel(executor utils.UserID)
	IsOpenToBattel(user utils.UserID) bool
	Battle(executor, rival utils.UserID) (bool, error)
	ModifyCards(owner utils.UserID, cards []CardWithUserInfluence) ([]CardParams, error)
}

type Battle struct {
	cardSetContract CardSetI
	cardContract    CardsI
	system          system.SystemI
	openToBattel    map[utils.UserID]bool
}

func NewBattle(
	cardSetContract CardSetI,
	cardContract CardsI,
	system system.SystemI,
) *Battle {
	return &Battle{
		cardSetContract: cardSetContract,
		cardContract:    cardContract,
		system:          system,
		openToBattel:    map[utils.UserID]bool{},
	}
}

func (b *Battle) ReadyToBattle(executor utils.UserID) {
	b.openToBattel[executor] = true
}

func (b *Battle) NotReadyToBattel(executor utils.UserID) {
	b.openToBattel[executor] = false
}

func (b *Battle) IsOpenToBattel(user utils.UserID) bool {
	open, ok := b.openToBattel[user]
	if !ok {
		open = false
	}
	return open
}

func (b *Battle) modifyCard(cardParam, userParam *CardParams, limit *system.UsersParam) {

	if limit.Kw >= userParam.Hp {
		limit.Kw -= userParam.Hp
		cardParam.Hp += userParam.Hp
	}

	if limit.Volts >= userParam.Strength {
		limit.Volts -= userParam.Strength
		cardParam.Strength += userParam.Strength
	}

	if limit.Amperes >= userParam.Accuracy {
		limit.Amperes -= userParam.Accuracy
		cardParam.Accuracy += userParam.Accuracy
	}

	if limit.Cyberlinks >= uint64(userParam.Level) {
		limit.Cyberlinks -= uint64(userParam.Level)
		cardParam.Level += userParam.Level
	}
}

func (b *Battle) ModifyCards(owner utils.UserID, cards []CardWithUserInfluence) ([]CardParams, error) {
	res := make([]CardParams, len(cards))
	userParam, err := b.system.GetUserParam(owner)
	if err != nil {
		return nil, err
	}

	for num, card := range cards {
		if b.cardContract.IsOwner(card.cardID, owner) != nil {
			continue
		}
		if b.cardContract.IsFreezed(card.cardID) {
			continue
		}

		cardParam, err := b.cardContract.GetCardProperties(card.cardID)
		if err != nil {
			continue
		}
		b.modifyCard(&cardParam, &card.userAttributes, &userParam)
		res[num] = cardParam
	}
	return res, nil
}

func (b *Battle) Battle(executor, rival utils.UserID) (bool, error) {
	// if !b.IsOpenToBattel(executor) {
	// 	return utils.ErrBattleUserIsNotReady
	// }

	if !b.IsOpenToBattel(rival) {
		return false, utils.ErrBattleUserIsNotReady
	}

	cardSetExecutor := b.cardSetContract.GetActualSetWithAttribute(executor)
	cardSetRival := b.cardSetContract.GetActualSetWithAttribute(rival)

	cardParamExecutor, err := b.ModifyCards(executor, cardSetExecutor)
	if err != nil {
		return false, err
	}
	cardParamRival, err := b.ModifyCards(rival, cardSetRival)
	if err != nil {
		return false, err
	}

	return b.battle(cardParamExecutor, cardParamRival), nil
}

func (b *Battle) battle(cardSet1, cardSet2 []CardParams) bool {
	sumCardsParam1 := CardParams{}
	sumCardsParam2 := CardParams{}
	for _, card := range cardSet1 {
		sumCardsParam1.Hp += card.Hp
		sumCardsParam1.Accuracy += card.Accuracy
		sumCardsParam1.Level += card.Level
		sumCardsParam1.Strength += card.Strength
	}
	for _, card := range cardSet1 {
		sumCardsParam2.Hp += card.Hp
		sumCardsParam2.Accuracy += card.Accuracy
		sumCardsParam2.Level += card.Level
		sumCardsParam2.Strength += card.Strength
	}

	steps1 := float64(sumCardsParam1.Hp) / float64(sumCardsParam2.Accuracy-sumCardsParam1.Strength)
	steps2 := float64(sumCardsParam2.Hp) / float64(sumCardsParam1.Accuracy-sumCardsParam2.Strength)

	return steps1 > steps2
}
