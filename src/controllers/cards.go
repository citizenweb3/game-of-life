package controllers

import (
	"encoding/json"

	"gameoflife/contracts"
	"gameoflife/utils"

	"net/http"

	log "github.com/sirupsen/logrus"
)

type CardsControllerI interface {
	Transfer(w http.ResponseWriter, r *http.Request)
	Burn(w http.ResponseWriter, r *http.Request)
	Freeze(w http.ResponseWriter, r *http.Request)
	Unfreeze(w http.ResponseWriter, r *http.Request)
	// AddAvatar(w http.ResponseWriter, r *http.Request)
	MintNewCard(w http.ResponseWriter, r *http.Request)

	GetCardProperties(w http.ResponseWriter, r *http.Request)
	GetOwnersCards(w http.ResponseWriter, r *http.Request)
}

type CardsController struct {
	card contracts.CardsI
}

func NewCardsController(card contracts.CardsI) *CardsController {
	return &CardsController{
		card: card,
	}
}

func (cc *CardsController) MintNewCard(w http.ResponseWriter, r *http.Request) {
	var p MintNewCardRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("userid is requred:")
		return
	}
	cardID, err := cc.card.MintNewCard(utils.UserID(p.UserID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("internal error:" + err.Error())
		return
	}

	params, err := cc.card.GetCardProperties(cardID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error: " + err.Error())
		return
	}

	resp := MintNewCardResponse{
		CardId:   cardID.ToString(),
		UserID:   p.UserID,
		Hp:       params.Hp,
		Level:    params.Level,
		Strength: params.Strength,
		Accuracy: params.Accuracy,
	}
	log.Info("resp ", resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (cc *CardsController) GetCardProperties(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	cardID := utils.CardID(q.Get("card_id"))
	if cardID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("card_id is required")
		return
	}
	prop, errProp := cc.card.GetCardProperties(cardID)
	owner, errOwner := cc.card.GetCardOwner(cardID)

	if errProp != nil || errOwner != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error " + errProp.Error() + errOwner.Error())
		return
	}
	resp := GetCardPropertiesResponce{
		CardId:   cardID.ToString(),
		Hp:       prop.Hp,
		Accuracy: prop.Accuracy,
		Level:    prop.Level,
		Strength: prop.Strength,
		UserID:   owner.ToString(),
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (cc *CardsController) GetOwnersCards(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	owner := q.Get("owner")
	if owner == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("owner is required")
		return
	}

	cards := cc.card.GetOwnersCards(utils.UserID(owner))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func (cc *CardsController) Transfer(w http.ResponseWriter, r *http.Request) {
	var p TransferRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.Executor == "" || p.To == "" || p.CardID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("executor, to, cardid are requered")
		return
	}

	err = cc.card.Transfer(utils.CardID(p.CardID), utils.UserID(p.Executor), utils.UserID(p.To))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Done")
}
func (cc *CardsController) Burn(w http.ResponseWriter, r *http.Request) {
	var p BurnRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.Executor == "" || p.CardID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("executor, cardid are requered")
		return
	}

	err = cc.card.Burn(utils.CardID(p.CardID), utils.UserID(p.Executor))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Done")
}

func (cc *CardsController) Freeze(w http.ResponseWriter, r *http.Request) {
	var p FreezeRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.Executor == "" || p.CardID1 == "" || p.CardID2 == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("executor, cardids are requered")
		return
	}

	err = cc.card.Freeze(utils.CardID(p.CardID1), utils.CardID(p.CardID2), utils.UserID(p.Executor))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Done")
}
func (cc *CardsController) Unfreeze(w http.ResponseWriter, r *http.Request) {

	var p BurnRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.Executor == "" || p.CardID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("executor, cardid are requered")
		return
	}

	err = cc.card.UnFreeze(utils.CardID(p.CardID), utils.UserID(p.Executor))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Done")
}
