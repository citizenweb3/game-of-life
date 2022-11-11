package controllers

import (
	"encoding/json"
	"gameoflife/contracts"
	"gameoflife/utils"
	"net/http"
)

type BattleControllerI interface {
	ReadyToBattle(w http.ResponseWriter, r *http.Request)
	IsOpenToBattel(w http.ResponseWriter, r *http.Request)
	Battle(w http.ResponseWriter, r *http.Request)
	ModifyCards(w http.ResponseWriter, r *http.Request)
}

type BattleController struct {
	battleContract contracts.BattleI
}

func NewBattleController(battleContract contracts.BattleI) *BattleController {
	return &BattleController{
		battleContract: battleContract,
	}
}

func (bc *BattleController) ReadyToBattle(w http.ResponseWriter, r *http.Request) {
	var p SetBattleReadinessRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.Executor == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("executor is requered")
		return
	}

	var outStr string
	if p.Ready {
		bc.battleContract.ReadyToBattle(utils.UserID(p.Executor))
		outStr = "Now your set is open for battle"
	} else {
		bc.battleContract.NotReadyToBattel(utils.UserID(p.Executor))
		outStr = "Now your set is protected from battle"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(outStr)

}
func (bc *BattleController) IsOpenToBattel(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := utils.CardID(q.Get("user_id"))
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("user_id is required")
		return
	}
	var resp string
	if bc.battleContract.IsOpenToBattel(utils.UserID(userID)) {
		resp = "This user is open to battle"
	} else {
		resp = "This user is not open to battle"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (bc *BattleController) Battle(w http.ResponseWriter, r *http.Request) {
	var p BattleRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.Executor == "" || p.Rival == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("executor and rival are requered")
		return
	}

	res, err := bc.battleContract.Battle(utils.UserID(p.Executor), utils.UserID(p.Rival))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Some problems. Choose peace. Enough fighting")
		return
	}
	var outStr string
	switch res {
	case -1:
		outStr = "Now are loose!!! But relax, I fogot add monitisation"
	case 1:
		outStr = "Now are win!!! And now that's all."
	case 0:
		outStr = "Nobody won. Fighting for too long."

	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(outStr)
}
func (bc *BattleController) ModifyCards(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Not now.")
}
