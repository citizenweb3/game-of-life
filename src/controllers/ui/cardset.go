package userinterface

import (
	"fmt"
	"gameoflife/utils"
	"html/template"
	"log"
	"net/http"
)

var cardSetPage = `<html><body> 
<head> ` + style + ` </head>
<body> ` + getMenu() + `
<h2>Card Set:</h2>

{{block "batch" .}}
<div>
	<div float="left" width="50%">
		<div class="vertical-menu">
		{{range .usersId}}
			<a href="cardsset?user_id={{.UserId}}">{{.UserId}}</a>
		{{end}}
		</div>
	</div>

	<div float="left"  width="50%">
		CardId to add|change <input list='card_ids' id="card_id"> <br>
		<button onclick="addCardToSet({{.userId}}, 0)">Add card to set to pos 0</button>
		<datalist id="card_ids">
			{{range .cardIds}}
				<option value="{{.CardId}}">
			{{end}}
		</datalist>
		<table>
		<tr>
			<td>CardId</td>
			<td>Hp</td>
			<td>Level</td>
			<td>Strength</td>
			<td>Accuracy</td>
			<td>Remove</td>
			<td>Change</td>
			<td>Add</td>
		</tr>
		{{range .cards}}
			<tr>
				<td>{{.CardId}}</td>
				<td>{{.Hp}}</td>
				<td>{{.Level}}</td>
				<td>{{.Strength}}</td>
				<td>{{.Accuracy}}</td>
				<td><button onclick="removeCardFromSet({{.UserId}}, {{.CardId}})">remove</button></td>
				<td><button onclick="changeCardFromSet({{.UserId}}, {{.CardId}})">Change</button></td>
				<td><button onclick="addCardToSet({{.UserId}}, {{.NumInSet}})">Add card to set</button> </td>
			</tr>
		{{end}}
		</table>
	</div>
{{end}}
<script>

function addCardToSet(userId, numInSet) {
	var cardId = document.getElementById('card_id').value
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/set/card", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "CardID": cardId, "NumInSet": numInSet}));} catch (err) { console.log(err) }
}
function removeCardFromSet(userId, cardId) {
	var xhr = new XMLHttpRequest();
	xhr.open("DELETE", "/set/card", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "CardID": cardId}));} catch (err) { console.log(err) }
}
function changeCardFromSet(userId, cardIdLast) {
	var cardIdNew = document.getElementById('card_id').value
	var xhr = new XMLHttpRequest();
	xhr.open("PATCH", "/set/card", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "CardIDLast": cardIdLast, "CardIDNew": cardIdNew}));} catch (err) { console.log(err) }
}
</script>
</body></html>`

type CardSetUsersSetsUI struct {
	UserId string
	Cards  []CardSetCardInfoUI
}
type CardSetCardInfoUI struct {
	UserId   string
	CardId   string
	Hp       string
	Level    string
	Strength string
	Accuracy string
	NumInSet int
}
type UserIDUI struct {
	UserId string
}

type CardIDUI struct {
	CardId string
}

func (ui *UI) GetCardSetPage(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	var userCardIds []CardIDUI
	for _, userCardId := range ui.cards.GetOwnersCards(utils.UserID(userId)) {
		userCardIds = append(userCardIds, CardIDUI{CardId: userCardId.Id.ToString()})
	}

	var userIds []UserIDUI
	for _, user := range ui.system.GetUserList() {
		userIds = append(userIds, UserIDUI{UserId: user.ToString()})
	}

	var cardsSet []CardSetCardInfoUI
	act := ui.cardSet.GetActualSet(utils.UserID(userId))
	for num, cardId := range act {
		prop, _ := ui.cards.GetCardProperties(cardId)
		cardsSet = append(cardsSet, CardSetCardInfoUI{
			UserId:   userId,
			CardId:   cardId.ToString(),
			Hp:       fmt.Sprint(prop.Hp),
			Accuracy: fmt.Sprint(prop.Accuracy),
			Strength: fmt.Sprint(prop.Strength),
			Level:    fmt.Sprint(prop.Level),
			NumInSet: num,
		})
	}
	data := map[string]interface{}{"userId": userId, "cards": cardsSet, "usersId": userIds, "cardIds": userCardIds}
	var t = template.Must(template.New("").Parse(cardSetPage))

	var err error
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
