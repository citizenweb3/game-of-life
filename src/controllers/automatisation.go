package controllers

import (
	"encoding/json"
	"fmt"
	"gameoflife/contracts"
	"gameoflife/system"
	"gameoflife/utils"
	"net/http"
)

type Automatisation struct {
	system   system.SystemI
	cards    contracts.CardsI
	cardsSet contracts.CardSetI
	battle   contracts.BattleI
}

func NewAutomatisation(
	system system.SystemI,
	cards contracts.CardsI,
	cardsSet contracts.CardSetI,
	battle contracts.BattleI,
) *Automatisation {
	return &Automatisation{
		system:   system,
		cards:    cards,
		cardsSet: cardsSet,
		battle:   battle,
	}
}

func (a *Automatisation) CreateSystem(w http.ResponseWriter, r *http.Request) {
	var p CreateSystemRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}

	if p.UserCountTo == 0 {
		if p.UserCountFrom == 0 {
			p.UserCountFrom = 1
		}
		p.UserCountTo = p.UserCountFrom + 3
	} else if p.UserCountFrom > p.UserCountTo {
		p.UserCountFrom, p.UserCountTo = p.UserCountTo, p.UserCountFrom
	}

	if p.CardCountTo == 0 {
		p.CardCountTo = p.CardCountFrom + 10
	} else if p.CardCountFrom > p.CardCountTo {
		p.CardCountFrom, p.CardCountTo = p.CardCountTo, p.CardCountFrom
	}

	var additionalCount int
	userCount := p.UserCountFrom
	if p.UserCountTo != p.UserCountFrom {
		additionalCount = int(utils.GetRandomNumberUint64(uint64(p.UserCountTo - p.UserCountFrom)))
		userCount += additionalCount
	}
	fmt.Printf("user created card count %d (from %d + add %d (to %d)) \n", userCount, p.UserCountFrom, p.UserCountTo, additionalCount)

	for userNum := 0; userNum < userCount; userNum++ {
		userName := fmt.Sprintf("random_user_%d_%d", utils.GetRandomNumberInt64(100000), userNum)
		userId := utils.UserID(userName)
		err := a.system.CreateUserWithRamdomParam(userId)

		if err != nil {
			continue
		}

		cardsCount := p.CardCountFrom
		if p.CardCountTo != p.CardCountFrom {
			cardsCount += int(utils.GetRandomNumberUint64(uint64(p.CardCountTo - p.CardCountFrom)))
		}

		fmt.Println("user created card count:", cardsCount)
		cardNumInSet := 0
		for cardNum := 0; cardNum < cardsCount; cardNum++ {
			cardId, err := a.cards.MintNewCard(userId)
			if err != nil {
				continue
			}
			if utils.GetRandomNumberBool() {
				a.cardsSet.AddCardToSet(userId, cardNumInSet, cardId)
				cardNumInSet++
			}
		}

		if utils.GetRandomNumberBool() {
			a.battle.ReadyToBattle(userId)
		}
	}
}
