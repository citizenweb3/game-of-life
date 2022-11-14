package userinterface

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var usersPage = `<html>
<head> ` + style + ` </head>
<body> ` + getMenu() + `
{{block "batch" .}}
	<div class="tab">
		<button class="tablinks" onclick="openTab(event, 'Users')" id="defaultOpen">Users</button>
		<button class="tablinks" onclick="openTab(event, 'SystemInfo')">SystemInfo</button>
		<button class="tablinks" onclick="openTab(event, 'GenerateRandomSystem')">GenerateRandomSystem</button>
	</div>
	<div id="Users" class="tabcontent">
		<h3>Users</h3>
		{{range .users}}
			<b><a href="cards?user_id={{.UserId}}"> {{.UserId}} </a></b> <br>
			<table>
			<tr><td>Volts</td><td>{{.Volts}}</td></tr>
			<tr><td>Amperes</td><td>{{.Amperes}}</td></tr>
			<tr><td>Cyberlinks</td><td>{{.Cyberlinks}}</td></tr>
			<tr><td>Kw</td><td>{{.Kw}}</td></tr>
			</table>
		{{end}}
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
	`
<script>
	function addUser() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/system/user/add", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"Random": true}));} catch (err) { console.log(err) }
	}
	function goToTheFuture() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/system/gotothefuture", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({"AddUnixTime": 2000}));} catch (err) { console.log(err) }
	}

	function getRandomSystem() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/generate", true);
		xhr.setRequestHeader("Content-Type", "application/json");
		try { xhr.send(JSON.stringify({}));} catch (err) { console.log(err) }
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

func (ui *UI) GetUsersPage(w http.ResponseWriter, r *http.Request) {
	// Put up some random data for demonstration:
	var userInfos []UserInfoUI
	usersIDs := ui.system.GetUserList()
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
	var t = template.Must(template.New("").Parse(usersPage))

	var err error
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func (ui *UI) GetUsersList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsersList")
	var userInfos []UserInfoUI
	usersIDs := ui.system.GetUserList()
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
	data := map[string]interface{}{"posts": userInfos}
	t := template.Must(template.New("").Parse(usersPage))

	err := t.ExecuteTemplate(w, "batch", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
