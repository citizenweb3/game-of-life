package controllers

import (
	"encoding/json"
	"fmt"
	"gameoflife/system"
	"gameoflife/utils"
	"math/rand"
	"net/http"
)

type SystemControllerI interface {
	GetUserList(w http.ResponseWriter, r *http.Request)
	GenerateUsers(w http.ResponseWriter, r *http.Request)
	GetUserInfo(w http.ResponseWriter, r *http.Request)
	MoveForward(w http.ResponseWriter, r *http.Request)
}

type SystemController struct {
	System *system.System
}

func NewSystemController(system *system.System) *SystemController {
	return &SystemController{
		System: system,
	}
}

func (sc *SystemController) GetUserList(w http.ResponseWriter, r *http.Request) {
	userList := sc.System.GetUserList()

	if len(userList) == 0 {
		return
	}
	resp := make(GetUserListResponse, 0, len(userList))

	for _, userID := range userList {
		userParam, err := sc.System.GetUserParam(userID)
		if err != nil {
			// todo add logging error
			continue
		}
		ui := UserInfoApi{
			UserID: string(userID),
		}
		ui.AddParams(userParam)
		resp = append(resp, ui)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (sc *SystemController) GenerateUsers(w http.ResponseWriter, r *http.Request) {
	var p UserInfoApi
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	if p.UserID == "" {
		p.UserID = getRandomUserID()
	}
	if p.Random == true {
		err = sc.System.CreateUserWithRamdomParam(utils.UserID(p.UserID))
	} else {
		err = sc.System.CreateUserWithParam(utils.UserID(p.UserID), p.GetParams())
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in creating:" + err.Error())
		return
	}

	userParam, err := sc.System.GetUserParam(utils.UserID(p.UserID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Can't find user:" + err.Error())
		return
	}
	resp := UserInfoApi{UserID: p.UserID}
	resp.AddParams(userParam)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (sc *SystemController) MoveForward(w http.ResponseWriter, r *http.Request) {
	var p MoveForwardRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("error in decode:" + err.Error())
		return
	}
	totalMoved := sc.System.MoveCurrentTimeForvard(p.AddUnixTime)

	resp := MoveForwardResponse{
		TotalMoved:  totalMoved,
		CurrentTime: sc.System.GetCurrentTime(),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (sc *SystemController) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := q.Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("user_id is required:")
		return
	}

	param, err := sc.System.GetUserParam(utils.UserID(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := UserInfoApi{UserID: userID}
	resp.AddParams(param)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(param)
}

func getRandomUserID() string {
	return fmt.Sprintf("random %d", rand.Int31n(10000000))
}
