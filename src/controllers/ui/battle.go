package userinterface

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

var battlePage = `<html><body> 
<head> ` + style + ` </head>
<body> ` + GetMenu("battle") + `
<h2>Battle</h2>

{{block "batch" .}}
<div>
	<div>
		<input list='users_id_list' id="user_id1"> vs <input list='users_id_list' id="user_id2"> <br>
		<datalist id="users_id_list">
			{{range .usersId}}
				<option value="{{.UserId}}">
			{{end}}
		</datalist>
		<button onclick="battle()">Battle!</button> 
	</div>
	
	<div>
	<input list='user_id_opened' id="user_id_to_close">
	<datalist id="user_id_opened">
		{{range .usersIdOpen}}
			<option value="{{.UserId}}">
		{{end}}
	</datalist>
	<button onclick="closeToBattle()">Close to battle</button>
	</div>

	<div>
	<input list='user_id_closed' id="user_id_to_open">
	<datalist id="user_id_closed">
		{{range .usersIdClose}}
			<option value="{{.UserId}}">
		{{end}}
	</datalist>
	<button onclick="openToBattle()">Open to battle</button>
	</div>

	<div>
	Open to Battle:
	<table>
		<tr>
			{{range .usersIdOpen}}
				<td>{{.UserId}}</td>
			{{end}}
		</tr>
	</table>
	</div>
{{end}}
<script>

function battle() {
	var userId1 = document.getElementById('user_id1').value
	var userId2 = document.getElementById('user_id2').value
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/battle/start", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId1, "Rival": userId2}));} catch (err) { console.log(err) }
}

function closeToBattle() {
	var userId = document.getElementById('user_id_to_close').value
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/battle/ready", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "Ready": false}));} catch (err) { console.log(err) }
	document.location.reload(true)
}
function openToBattle() {
	var userId = document.getElementById('user_id_to_open').value
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/battle/ready", true);
	xhr.setRequestHeader("Content-Type", "application/json");
	try { xhr.send(JSON.stringify({"Executor": userId, "Ready": true}));} catch (err) { console.log(err) }
	document.location.reload(true)
}
</script>
</body></html>`

func (ui *UI) GetBattlePage(w http.ResponseWriter, r *http.Request) {

	var userIds []UserIDUI
	var openToBattle []UserIDUI
	var notOpenToBattle []UserIDUI

	users := ui.system.GetUserList()
	sort.Slice(users, func(i, j int) bool { return -1 == strings.Compare(users[i].ToString(), users[j].ToString()) })

	for _, user := range users {
		userIds = append(userIds, UserIDUI{UserId: user.ToString()})
		if ui.battle.IsOpenToBattel(user) {
			openToBattle = append(openToBattle, UserIDUI{UserId: user.ToString()})
		} else {
			notOpenToBattle = append(notOpenToBattle, UserIDUI{UserId: user.ToString()})
		}
	}

	fmt.Println("userIds", userIds)
	fmt.Println("open", openToBattle)
	fmt.Println("close", notOpenToBattle)

	data := map[string]interface{}{"usersId": userIds, "usersIdOpen": openToBattle, "usersIdClose": notOpenToBattle}
	var t = template.Must(template.New("").Parse(battlePage))

	var err error
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
