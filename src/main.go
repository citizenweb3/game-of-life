package main

import (
	"gameoflife/api"
	"gameoflife/contracts"
	"gameoflife/controllers"
	user_interface "gameoflife/controllers/ui"
	"gameoflife/system"
	"gameoflife/utils/config"
)

func main() {
	system := system.NewSystem()
	systemController := controllers.NewSystemController(system)

	card := contracts.NewCards(system, 1000)
	cardsController := controllers.NewCardsController(card)

	cardSet := contracts.NewCardSet(card, uint8(5))
	cardSetController := controllers.NewCardSetController(cardSet, card)

	battle := contracts.NewBattle(cardSet, card, system)
	battleController := controllers.NewBattleController(battle)

	au := controllers.NewAutomatisation(system, card, cardSet, battle)

	cfg := config.NewConfigApp()

	ui := user_interface.NewUI(system, card, cardSet, battle)

	httpServer := api.NewHttpServer(cfg, systemController, cardsController, cardSetController, battleController, au, ui)

	defer httpServer.Stop()
	httpServer.Start()
}
