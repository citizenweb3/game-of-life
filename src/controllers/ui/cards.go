package userinterface

import (
	"fmt"
	"gameoflife/utils"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

var cardsPage = `<html>
<head> ` + style + ` </head>
<body> ` + GetMenu("cards") + `

{{block "batch" .}}
<h2>User {{.UserId}} cards</h2>

<div>
	<div float="left" width="50%">
		<div class="horisontal-menu">
			{{range .usersId}}
				<a href="cards?user_id={{.UserId}}">{{.UserId}}</a>
			{{end}}
		</div>
	</div>
	<div float="left" width="50%">
		<button onclick="mintNewCard()">Mint card</button> 

		<div>
			<button onclick="transferCard()">transfer</button> 
			CardId <input type="text" id = "card_id" placeholder ="card_id">
			Receiver <input type="text" id = "receiver" placeholder ="receiver">
		</div>

		<div>
			<button onclick="freezeCard()">Freeze</button> 
			CardId <input type="text" id="card_id_freeze1" placeholder ="card_id1">
			CardId <input type="text" id="card_id_freeze2" placeholder ="card_id2">
		</div>
		<div>			
			<button onclick="unfreezeCard()">UnFreeze</button> 
			CardId <input type="text" id="card_id_unfreeze" placeholder ="card_id">
		</div>

		<div id="cards">
			<table>
				<tr>
					<td>CardId</td>
					<td>Hp</td>
					<td>Level</td>
					<td>Deffence</td>
					<td>Expected Damage<br>(Accuracy-Damage)</td>
					<td>Freeze</td>
					<td>Burn</td>
				</tr>
			{{range .card_info}}
				<tr>
					<td>{{.CardIdShort}}</td>
					<td>{{.Hp}}</td>
					<td>{{.Level}}</td>
					<td>{{.Deffence}}</td>
					<td>{{.Accuracy}} - {{.Damage}}</td>
					<td></td>
					<td>{{.Freeze}}</td>
					<td><button onclick="burnCard('{{.CardId}}')">burn</button></td>
				</tr>
			{{end}}
			</table>
		</div>
	</div>
</div>
 ` + scripts + `
<script>

	function mintNewCard() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/cards/mint", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"UserID":{{.UserId}} }));} catch (err) { console.log(err) }
	}
	
	function transferCard() {
		var cardId = document.getElementById('card_id').value
		var receiver = document.getElementById('receiver').value
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/cards/transfer", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"Executor":{{.UserId}}, "CardID":cardId, "To":receiver}));} catch (err) { console.log(err) }
	}
	function burnCard() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/cards/burn", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"Executor":{{.UserId}}, "CardID":cardId}));} catch (err) { console.log(err) }
	}
	function freezeCard() {		
		var cardId1 = document.getElementById('card_id_freeze1').value
		var cardId2 = document.getElementById('card_id_freeze2').value
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/cards/freeze", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"Executor":{{.UserId}}, "CardID1":cardId1, "CardID2":cardId2}));} catch (err) { console.log(err) }
	}
	function unfreezeCard() {
		var cardId = document.getElementById('card_id_unfreeze').value
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/cards/unfreeze", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"Executor":{{.UserId}}, "CardID":cardId}));} catch (err) { console.log(err) }
	}
</script>
{{end}}
</body></html>`

type CardsInfoUI struct {
	CardId      string
	CardIdShort string
	Hp          string
	Level       string
	Deffence    string
	Damage      string
	Accuracy    string
	Freeze      string
}

func (ui *UI) GetCardsPage(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	users := ui.system.GetUserList()
	sort.Slice(users, func(i, j int) bool { return -1 == strings.Compare(users[i].ToString(), users[j].ToString()) })

	var userIds []UserIDUI
	for _, user := range users {
		userIds = append(userIds, UserIDUI{UserId: user.ToString()})
	}
	var userInfos []CardsInfoUI
	cards := ui.cards.GetOwnersCards(utils.UserID(userId))
	currTime := ui.system.GetCurrentTime()
	for _, card := range cards {
		freezeTime := ui.cards.GetFreezeTime(card.Id)
		freezeTimeStr := ""
		if freezeTime != 0 {
			freezeTimeStr = fmt.Sprint(freezeTime, "(after ", currTime-freezeTime, ")")
		}
		cardIdStr := card.Id.ToString()
		userInfos = append(userInfos, CardsInfoUI{
			CardId:      cardIdStr,
			CardIdShort: fmt.Sprintf("%s...%s", cardIdStr[:4], cardIdStr[len(cardIdStr)-4:]),
			Hp:          fmt.Sprint(card.Params.Hp),
			Level:       fmt.Sprint(card.Params.Level),
			Deffence:    fmt.Sprint(card.Params.Deffence),
			Accuracy:    fmt.Sprint(card.Params.Accuracy),
			Damage:      fmt.Sprint(card.Params.Damage),
			Freeze:      freezeTimeStr,
		})
	}
	data := map[string]interface{}{"card_info": userInfos, "UserId": userId, "usersId": userIds}
	var t = template.Must(template.New("").Parse(cardsPage))

	var err error
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// func (ui *UI) GetUsersList(w http.ResponseWriter, r *http.Request) {
// 	var userInfos []UserInfoUI
// 	usersIDs := ui.system.GetUserList()
// 	for _, userId := range usersIDs {
// 		param, err := ui.system.GetUserParam(userId)
// 		if err != nil {
// 			continue
// 		}
// 		userInfos = append(userInfos, UserInfoUI{
// 			UserId:     userId.ToString(),
// 			Amperes:    fmt.Sprint(param.GetAmperes()),
// 			Volts:      fmt.Sprint(param.GetVolts()),
// 			Kw:         fmt.Sprint(param.GetKw()),
// 			Cyberlinks: fmt.Sprint(param.GetCountCyberlinks()),
// 		})
// 	}
// 	data := map[string]interface{}{"posts": userInfos}
// 	t := template.Must(template.New("").Parse(usersPage))

// 	err := t.ExecuteTemplate(w, "batch", data)
// 	if err != nil {
// 		log.Printf("Template execution error: %v", err)
// 	}
// }
