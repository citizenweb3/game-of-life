package userinterface

import (
	"fmt"
	"gameoflife/contracts"
	"gameoflife/utils"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

var cardSetPage = `<html><body> 
<head> ` + style + ` </head>
<body> ` + GetMenu("cardsset") + `
<h2>Card Set:</h2>

{{block "batch" .}}
<div>
	<div float="left" width="50%">
		<div class="horisontal-menu">
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
			<td>Deffence</td>
			<td>Damage</td>
			<td>Remove</td>
			<td>Change</td>
			<td>Add</td>
		</tr>
		{{range .cards}}
			<tr>
				<td>{{.CardIdShort}}</td>
				<td>{{.Hp}}</td>
				<td>{{.Level}}</td>
				<td>{{.Deffence}}</td>
				<td>{{.Accuracy}} - {{.Damage}}</td>
				<td><button onclick="removeCardFromSet({{.UserId}}, {{.CardId}})">remove</button></td>
				<td><button onclick="changeCardFromSet({{.UserId}}, {{.CardId}})">Change</button></td>
				<td><button onclick="addCardToSet({{.UserId}}, {{.NumInSet}})">Add card to set</button> </td>
			</tr>
		{{end}}
		</table>
	</div>
	<h3>Influence:</h3>
	Set persent of user value to influence cardParam (0-100% of user value) <br>
	Num in set <input type="number" id="numInSet" name="numInSet" value=0 width="20px"><br>
	Hp <input type="number" id="hp" name="hp" value=0 width="20px"><br>
	Deffence <input type="number" id="def" name="def" value=0 width="20px"><br>
	Damage <input type="number" id="dam" name="dam" value=0 width="20px"><br>
	Accuracy <input type="number" id="acc" name="acc" value=0 width="20px"><br>
	
	<button onclick="addInfluence({{.userId}})">add influence</button>
{{end}}
<script>

function addInfluence(userId) {
	var numInSet = parseInt(document.getElementById('numInSet').value);
	var hp = parseInt(document.getElementById('hp').value);
	var def = parseInt(document.getElementById('def').value);
	var dam = parseInt(document.getElementById('dam').value);
	var acc = parseInt(document.getElementById('acc').value);

	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/set/attribute", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	var jsonVar = {"Executor": userId, "NumInSet": numInSet, "Hp": hp, "Deffence": def,  "Damage": dam,  "Accuracy": acc}
	try { xhr.send(JSON.stringify(jsonVar));} catch (err) { console.log(err) }
	document.location.reload(true)
}

function addCardToSet(userId, numInSet) {
	var cardId = document.getElementById('card_id').value
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/set/card", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "CardID": cardId, "NumInSet": numInSet}));} catch (err) { console.log(err) }
	document.location.reload(true)
}
function removeCardFromSet(userId, cardId) {
	var xhr = new XMLHttpRequest();
	xhr.open("DELETE", "/set/card", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "CardID": cardId}));} catch (err) { console.log(err) }
	document.location.reload(true)
}
function changeCardFromSet(userId, cardIdLast) {
	var cardIdNew = document.getElementById('card_id').value
	var xhr = new XMLHttpRequest();
	xhr.open("PATCH", "/set/card", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "CardIDLast": cardIdLast, "CardIDNew": cardIdNew}));} catch (err) { console.log(err) }
	document.location.reload(true)
}
</script>
</body></html>`

type CardSetUsersSetsUI struct {
	UserId string
	Cards  []CardSetCardInfoUI
}
type CardSetCardInfoUI struct {
	UserId      string
	CardId      string
	CardIdShort string
	Hp          string
	Level       string
	Deffence    string
	Damage      string
	Accuracy    string
	NumInSet    int
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
	usersIDs := ui.system.GetUserList()
	sort.Slice(usersIDs, func(i, j int) bool { return -1 == strings.Compare(usersIDs[i].ToString(), usersIDs[j].ToString()) })

	for _, user := range usersIDs {
		userIds = append(userIds, UserIDUI{UserId: user.ToString()})
	}

	var cardsSet []CardSetCardInfoUI
	act := ui.cardSet.GetActualSet(utils.UserID(userId))

	influenses := ui.cardSet.GetActualSetWithAttribute(utils.UserID(userId))
	for num, cardId := range act {
		infl := contracts.CardWithUserInfluence{}
		if len(influenses) > num {
			infl = influenses[num]
		}
		prop, _ := ui.cards.GetCardProperties(cardId)
		cardIdStr := cardId.ToString()
		cardIdShortStr := ""
		if len(cardIdStr) != 0 {
			cardIdShortStr = fmt.Sprintf("%s...%s,", cardIdStr[:4], cardIdStr[len(cardIdStr)-4:])
		}
		cardsSet = append(cardsSet, CardSetCardInfoUI{
			UserId:      userId,
			CardId:      cardIdStr,
			CardIdShort: cardIdShortStr,
			Hp:          fmt.Sprintf("%d (+%% %.1f)", prop.Hp, infl.UserAttributes.Hp),
			Accuracy:    fmt.Sprintf("%d (+%% %.1f)", prop.Accuracy, infl.UserAttributes.Accuracy),
			Deffence:    fmt.Sprintf("%d (+%% %.1f)", prop.Deffence, infl.UserAttributes.Deffence),
			Damage:      fmt.Sprintf("%d (+%% %.1f)", prop.Damage, infl.UserAttributes.Damage),
			Level:       fmt.Sprint(prop.Level),
			NumInSet:    num,
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
