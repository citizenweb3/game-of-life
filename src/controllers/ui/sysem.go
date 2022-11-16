package userinterface

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

var systemPage = `<html>
<head> ` + style + ` </head>
<body> ` + GetMenu("system") + `
{{block "batch" .}}
	<div class="tab">
		<button class="tablinks" onclick="openTab(event, 'Users')" id="defaultOpen">Users</button>
		<button class="tablinks" onclick="openTab(event, 'SystemInfo')">SystemInfo</button>
		<button class="tablinks" onclick="openTab(event, 'GenerateRandomSystem')">GenerateRandomSystem</button>
	</div>
	<div id="Users" class="tabcontent">
		<h3>Users</h3>
		Add value <input type="number" id="addvalue" name="addvalue" value=10>
		<div>
		{{range .users}}
			<div>
				<b><a href="cards?user_id={{.UserId}}"> {{.UserId}} </a></b> <br>
				<table>
				<tr><td>Volts     </td><td>{{.Volts}}     </td><td><button onclick='addParam({{.UserId}}, "Volts")'>Add</button></td></tr>
				<tr><td>Amperes   </td><td>{{.Amperes}}   </td><td><button onclick='addParam({{.UserId}}, "Amperes")'>Add</button></td></tr>
				<tr><td>Cyberlinks</td><td>{{.Cyberlinks}}</td><td><button onclick='addParam({{.UserId}}, "Cyberlinks")'>Add</button></td></tr>
				<tr><td>Kw        </td><td>{{.Kw}}        </td><td><button onclick='addParam({{.UserId}}, "Kw")'>Add</button></td></tr>
				</table>
			</div>
		{{end}}
		</div>
	</div>
	<div id="SystemInfo" class="tabcontent">
	<h3>SystemInfo</h3>
	<p> currentSystemTime: {{.Time}} </p>
	
	<button onclick="goToTheFuture()">go to the future 2000</button> 
	</div>
	<div id="GenerateRandomSystem" class="tabcontent">
		<h3>GenerateRandomSystem</h3>
		<p> Refrash page after generate</p>
		<button onclick="getRandomSystem()">Generate random system</button> 
		<button onclick="addUser()">Add user</button> 
	</div>
{{end}} ` + scripts +
	`<script>
	function addParam(userId, typeValueAdd) {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/system/user/param", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		var volts = 0;
		var amper = 0;
		var cyberlink = 0;
		var kw = 0;

		if (typeValueAdd == "Amperes" ){
			amper = parseInt(document.getElementById('addvalue').value)
		} else if (typeValueAdd == "Volts") {
			volts = parseInt(document.getElementById('addvalue').value)
		} else if (typeValueAdd == "Cyberlinks") {
			cyberlink = parseInt(document.getElementById('addvalue').value)
		} else if (typeValueAdd == "Kw" ){
			kw = parseInt(document.getElementById('addvalue').value)
		}

		var jsonVal = {"UserId": userId, "Volts": volts, "Amperes": amper, "Cyberlinks": cyberlink, "Kw": kw}
		try { xhr.send(JSON.stringify(jsonVal));} catch (err) { console.log(err) }
		document.location.reload(true)
	}
	function addUser() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/system/user/add", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"Random": true}));} catch (err) { console.log(err) }
		document.location.reload(true)
	}
	function goToTheFuture() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/system/gotothefuture", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"AddUnixTime": 2000}));} catch (err) { console.log(err) }
		document.location.reload(true)
	}

	function getRandomSystem() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/generate", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({}));} catch (err) { console.log(err) }
		document.location.reload(true)
	}
</script>
</body></html>
`

// JSON.stringify({"UserCountFrom": 10, "UserCountTo": 11, "CardCountFrom": 2, "CardCountTo": 20 })

type UserInfoUI struct {
	UserId     string
	Volts      string
	Amperes    string
	Cyberlinks string
	Kw         string
}

func (ui *UI) GetSystemPage(w http.ResponseWriter, r *http.Request) {
	// Put up some random data for demonstration:
	var userInfos []UserInfoUI
	usersIDs := ui.system.GetUserList()
	sort.Slice(usersIDs, func(i, j int) bool { return -1 == strings.Compare(usersIDs[i].ToString(), usersIDs[j].ToString()) })

	for _, userId := range usersIDs {
		param, err := ui.system.GetUserParam(userId)
		if err != nil {
			continue
		}
		userInfos = append(userInfos, UserInfoUI{
			UserId:     userId.ToString(),
			Amperes:    fmt.Sprint(param.GetAmperes()),
			Volts:      fmt.Sprint(param.GetVolts()),
			Kw:         fmt.Sprint(param.GetKw()),
			Cyberlinks: fmt.Sprint(param.GetCountCyberlinks()),
		})
	}
	currentTime := fmt.Sprint(ui.system.GetCurrentTime())
	data := map[string]interface{}{"users": userInfos, "Time": currentTime}
	var t = template.Must(template.New("").Parse(systemPage))

	var err error
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
