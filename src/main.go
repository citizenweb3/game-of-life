package main

import (
	"gameoflife/api"
	"gameoflife/contracts"
	"gameoflife/controllers"
	"gameoflife/system"
	"gameoflife/utils/config"
)

func main() {
	system := system.NewSystem()
	systemController := controllers.NewSystemController(system)

	card := contracts.NewCards(1000)
	cardsController := controllers.NewCardsController(card)

	cardSet := contracts.NewCardSet(card, uint8(100))
	cardSetController := controllers.NewCardSetController(cardSet, card)

	battle := contracts.NewBattle(cardSet, card, system)
	battleController := controllers.NewBattleController(battle)

	cfg := config.NewConfigApp()
	httpServer := api.NewHttpServer(cfg, systemController, cardsController, cardSetController, battleController)

	defer httpServer.Stop()
	httpServer.Start()
}
