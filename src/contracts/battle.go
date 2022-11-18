package contracts

import (
	"fmt"
	"gameoflife/system"
	"gameoflife/utils"
)

type BattleI interface {
	ReadyToBattle(executor utils.UserID)
	NotReadyToBattel(executor utils.UserID)
	IsOpenToBattel(user utils.UserID) bool
	Battle(executor, rival utils.UserID) (int, error)
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

func (b *Battle) modifyCard(cardParam *CardParams, userInfluene Influence, userParam CardParams) {
	cardParam.Hp += uint64(float64(userParam.Hp) * userInfluene.Hp)
	// cardParam.Level
	cardParam.Damage += uint64(float64(userParam.Damage) * userInfluene.Damage)
	cardParam.Deffence += uint64(float64(userParam.Deffence) * userInfluene.Deffence)
	cardParam.Accuracy += uint64(float64(userParam.Accuracy) * userInfluene.Accuracy)
}

func (b *Battle) ModifyCards(owner utils.UserID, cards []CardWithUserInfluence) ([]CardParams, error) {
	res := make([]CardParams, len(cards))
	userParam, err := b.system.GetUserParam(owner)
	if err != nil {
		return nil, err
	}
	userCardParam := CardParams{}
	userCardParam.ByUserParam(userParam)

	for num, card := range cards {
		cardParam, err := b.cardContract.GetCardProperties(card.CardID)
		if err != nil {
			continue
		}
		b.modifyCard(&cardParam, card.UserAttributes, userCardParam)
		res[num] = cardParam
	}
	return res, nil
}

func (b *Battle) Battle(executor, rival utils.UserID) (int, error) {
	// if !b.IsOpenToBattel(executor) {
	// 	return utils.ErrBattleUserIsNotReady
	// }

	if !b.IsOpenToBattel(rival) {
		return 0, utils.ErrBattleUserIsNotReady
	}

	cardSetExecutor := b.cardSetContract.GetActualSetWithAttribute(executor)
	if len(cardSetExecutor) == 0 {
		return 0, utils.ErrBattleUserHasNoOneCardsInSet
	}
	cardSetRival := b.cardSetContract.GetActualSetWithAttribute(rival)
	if len(cardSetRival) == 0 {
		return 0, utils.ErrBattleUserHasNoOneCardsInSet
	}
	cardParamExecutor, err := b.ModifyCards(executor, cardSetExecutor)
	if err != nil {
		return 0, err
	}
	cardParamRival, err := b.ModifyCards(rival, cardSetRival)
	if err != nil {
		return 0, err
	}
	// return b.battle(cardParamExecutor, cardParamRival), nil

	return b.battleStepByStep(cardParamExecutor, cardParamRival, 20), nil

}

func randDammage(deffence, damage, accuracy uint64) uint64 {
	if accuracy < damage {
		diff := damage - accuracy
		damage = damage + utils.GetRandomNumberUint64(diff) + accuracy
	}

	if deffence > damage {
		return 0
	}
	return damage - deffence
}

func (b *Battle) battleStep(cardSet1, cardSet2 []CardParams) bool {
	listCardAttack := make([]int, 0)
	allYourCardsDie := true
	for numInSet := 0; numInSet < len(cardSet1); numInSet++ {
		if cardSet1[numInSet].Hp != 0 {
			listCardAttack = append(listCardAttack, numInSet)
			allYourCardsDie = false
		}
		if len(listCardAttack) == 0 {
			continue
		}
		if cardSet2[numInSet].Hp == 0 {
			continue
		}

		sumAttack := uint64(0)
		for _, numCardSet1 := range listCardAttack {
			sumAttack += randDammage(cardSet2[numInSet].Deffence, cardSet1[numInSet].Damage, cardSet1[numCardSet1].Accuracy)
		}
		if cardSet2[numInSet].Hp < sumAttack {
			cardSet2[numInSet].Hp = 0
		} else {
			cardSet2[numInSet].Hp -= sumAttack
		}
	}
	return allYourCardsDie
}

func (b *Battle) battleStepByStep(cardSet1, cardSet2 []CardParams, maxStep int) int {
	fmt.Println("Start battle modified cardSet1=", cardSet1, " modified cardSet2=", cardSet2, "maxStep", maxStep)
	for i := 0; i < maxStep; i++ {
		if b.battleStep(cardSet1, cardSet2) {
			return -1
		}
		fmt.Println("Step ", i, " after executor hit rival cardSet1=", cardSet1, " cardSet2=", cardSet2)
		if b.battleStep(cardSet2, cardSet1) {
			return 1
		}
		fmt.Println("Step ", i, " after rival hit executor cardSet1=", cardSet1, " cardSet2=", cardSet2)
	}
	return 0
}
