package api

import (
	"fmt"
	"gameoflife/controllers"
	"gameoflife/utils/config"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func InitializeRoutes(
	systemController controllers.SystemControllerI,
	cardsController controllers.CardsControllerI,
	cardSetController controllers.CardSetControllerI,
	battleController controllers.BattleControllerI,
) *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	// system
	router.HandleFunc("/system/user/list", systemController.GetUserList).Methods("GET")
	router.HandleFunc("/system/user/add", systemController.GenerateUsers).Methods("POST")
	router.HandleFunc("/system/user", systemController.GetUserInfo).Methods("GET")
	router.HandleFunc("/system/gotothefuture", systemController.MoveForward).Methods("GET")

	// Cards
	router.HandleFunc("/cards/mint", cardsController.MintNewCard).Methods("POST")
	router.HandleFunc("/cards/list", cardsController.GetOwnersCards).Methods("GET")
	router.HandleFunc("/cards/info", cardsController.GetCardProperties).Methods("GET")
	router.HandleFunc("/cards/transfer", cardsController.Transfer).Methods("POST")
	router.HandleFunc("/cards/burn", cardsController.Burn).Methods("POST")
	router.HandleFunc("/cards/freeze", cardsController.Freeze).Methods("POST")
	router.HandleFunc("/cards/unfreeze", cardsController.Unfreeze).Methods("POST")

	// CardsSet
	router.HandleFunc("/set/actual", cardSetController.GetActualSet).Methods("GET")
	router.HandleFunc("/set/card", cardSetController.AddCardToSet).Methods("POST")
	router.HandleFunc("/set/card", cardSetController.ChangeCardFromSet).Methods("PATCH")
	router.HandleFunc("/set/card", cardSetController.RemoveCardFromSet).Methods("DELETE")

	router.HandleFunc("/set/attribute", cardSetController.SetUserAttribute).Methods("POST")
	router.HandleFunc("/set/attribute", cardSetController.GetUserAttributes).Methods("GET")

	// Battle
	router.HandleFunc("/battle/start", battleController.Battle).Methods("POST")
	router.HandleFunc("/battle/isopen", battleController.IsOpenToBattel).Methods("GET")
	router.HandleFunc("/battle/ready", battleController.ReadyToBattle).Methods("POST")

	return router

}

type HttpServer struct {
	cfg    *config.ConfigApp
	router *mux.Router
	server *http.Server
}

func NewHttpServer(
	cfg *config.ConfigApp,
	sc controllers.SystemControllerI,
	cc controllers.CardsControllerI,
	csc controllers.CardSetControllerI,
	bc controllers.BattleControllerI,
) *HttpServer {
	return &HttpServer{
		cfg:    cfg,
		router: InitializeRoutes(sc, cc, csc, bc),
		server: nil,
	}
}

func (server *HttpServer) Start() error {
	// Start Server
	server.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", server.cfg.Port),
		Handler: server.router,
	}

	fmt.Println("HttpServer Start")
	return server.server.ListenAndServe()
}

func (server *HttpServer) Stop() {
	err := server.server.Close()
	if err != nil {
		log.Info("error when stopped http server %s", err)
	}
}
