package userinterface

import "fmt"

type MenuInfo struct {
	Link     string
	Text     string
	IsActive bool
}

func GetMenu(actualLink string) string {
	var menuInfo = []MenuInfo{
		{Link: "system", Text: "System", IsActive: "system" == actualLink},
		{Link: "users", Text: "Users", IsActive: "users" == actualLink},
		{Link: "cards", Text: "Cards", IsActive: "cards" == actualLink},
		{Link: "cardsset", Text: "Cards set", IsActive: "cardsset" == actualLink},
		{Link: "battle", Text: "Battle", IsActive: "battle" == actualLink},
	}
	return getMenu(menuInfo)
}

func getMenu(links []MenuInfo) string {
	menu := `<div class="vertical-menu">`
	for _, menuInfo := range links {
		currentStyle := ""
		if menuInfo.IsActive {
			currentStyle = `style="background-color: #ccf;"`
		}
		menu += fmt.Sprintf("<a href='%s' %s> %s </a>\n", menuInfo.Link, currentStyle, menuInfo.Text)
	}

	return menu + "</div>\n<br><br>"
}
