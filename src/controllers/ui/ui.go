package userinterface

import (
	"gameoflife/system"

	"gameoflife/contracts"
)

type UI struct {
	system  system.SystemI
	cards   contracts.CardsI
	cardSet contracts.CardSetI
	battle  contracts.BattleI
}

func NewUI(system system.SystemI, cards contracts.CardsI, cardSet contracts.CardSetI, battle contracts.BattleI) *UI {
	return &UI{
		system:  system,
		cards:   cards,
		cardSet: cardSet,
		battle:  battle,
	}
}
