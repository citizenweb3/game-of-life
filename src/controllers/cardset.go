package controllers

import (
	"encoding/json"
	"fmt"
	"gameoflife/contracts"
	"gameoflife/utils"
	"net/http"
)

type CardSetControllerI interface {
	GetActualSet(w http.ResponseWriter, r *http.Request)
	AddCardToSet(w http.ResponseWriter, r *http.Request)
	RemoveCardFromSet(w http.ResponseWriter, r *http.Request)
	ChangeCardFromSet(w http.ResponseWriter, r *http.Request)

	SetUserAttribute(w http.ResponseWriter, r *http.Request)
	GetUserAttributes(w http.ResponseWriter, r *http.Request)
}

type CardSetController struct {
	cardSetContract contracts.CardSetI
	cardsContract   contracts.CardsI
}

func NewCardSetController(
	cardSetContract contracts.CardSetI,
	cardsContract contracts.CardsI,
) *CardSetController {
	return &CardSetController{
		cardSetContract: cardSetContract,
		cardsContract:   cardsContract,
	}
}

func (csc *CardSetController) GetActualSet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := utils.CardID(q.Get("user_id"))
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("user_id is required")
		return
	}

	cardIDs := csc.cardSetContract.GetActualSet(utils.UserID(userID))
	resp := make(GetActualSetResponse, 0, len(cardIDs))
	for _, card := range cardIDs {
		prop, errpr := csc.cardsContract.GetCardProperties(card)
		owner, errow := csc.cardsContract.GetCardOwner(card)
		if errpr != nil {
			prop = contracts.CardParams{}
		}
		if errow != nil {
			owner = ""
		}
		resp = append(resp, GetActualSetResponseElement{
			CardID:   string(card),
			Owner:    string(owner),
			Hp:       prop.Hp,
			Level:    prop.Level,
			Deffence: prop.Deffence,
			Damage:   prop.Damage,
			Accuracy: prop.Accuracy,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (csc *CardSetController) AddCardToSet(w http.ResponseWriter, r *http.Request) {
	var p AddCardToSetRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		responseWithErrorDecode(w, err)
		return
	}

	if p.CardID == "" || p.Executor == "" || p.NumInSet < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("card id and execuer are required")
		return
	}

	err = csc.cardSetContract.AddCardToSet(utils.UserID(p.Executor), p.NumInSet, utils.CardID(p.CardID))

	if err != nil {
		responseWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("done")
}

func (csc *CardSetController) RemoveCardFromSet(w http.ResponseWriter, r *http.Request) {
	var p RemoveCardToSetRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		responseWithErrorDecode(w, err)
		return
	}

	if p.CardID == "" || p.Executor == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("cardid and execuer are required")
		return
	}

	err = csc.cardSetContract.RemoveCardFromSet(utils.UserID(p.Executor), utils.CardID(p.CardID))

	if err != nil {
		responseWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("done")
}

func (csc *CardSetController) ChangeCardFromSet(w http.ResponseWriter, r *http.Request) {
	var p ChangeCardFromSetRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		responseWithErrorDecode(w, err)
		return
	}

	if p.CardIDLast == "" || p.CardIDNew == "" || p.Executor == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("cardidlast, cardidnew and execuer are required")
		return
	}

	err = csc.cardSetContract.ChangeCardFromSet(utils.UserID(p.Executor), utils.CardID(p.CardIDLast), utils.CardID(p.CardIDNew))

	if err != nil {
		responseWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("done")
}

func (csc *CardSetController) SetUserAttribute(w http.ResponseWriter, r *http.Request) {

	fmt.Println("SetUserAttribute")
	var p SetUserAttributeRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		responseWithErrorDecode(w, err)
		return
	}
	fmt.Println("SetUserAttribute", p)

	err = csc.cardSetContract.SetUserAttribute(utils.UserID(p.Executor), uint8(p.NumInSet),
		contracts.Influence{
			Hp:       p.Hp,
			Deffence: p.Deffence,
			Damage:   p.Damage,
			Accuracy: p.Accuracy,
		})

	if err != nil {
		responseWithError(w, err)
		return
	}

	fmt.Println("SetUserAttribute Done")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("done")
}

func (csc *CardSetController) GetUserAttributes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := utils.CardID(q.Get("user_id"))
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("user_id is required")
		return
	}

	userAttr := csc.cardSetContract.GetUserAttributes(utils.UserID(userID))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userAttr)
}
