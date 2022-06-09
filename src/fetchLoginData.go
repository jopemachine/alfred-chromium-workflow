package src

import (
	"fmt"

	"github.com/deanishe/awgo"
)

var FetchLoginData = func(wf *aw.Workflow, query string) {
	var dbQuery = fmt.Sprintf(`
		SELECT username_element, username_value, origin_url
			FROM logins
			WHERE origin_url LIKE '%%%s%%' OR username_element LIKE '%%%s%%' OR username_value LIKE '%%%s%%'
			ORDER BY date_last_used
	`, query, query, query)

	loginDataDB := GetLoginDataDB(wf)

	rows, err := loginDataDB.Query(dbQuery)
	CheckError(err)

	var userNameValue string
	var userNameElement string
	var originUrl string

	for rows.Next() {
		err := rows.Scan(&userNameElement, &userNameValue, &originUrl)
		CheckError(err)

		if userNameValue == "" {
			continue
		}

		domainName := ExtractDomainName(originUrl)
		iconPath := fmt.Sprintf(`%s/%s.png`, GetFaviconDirectoryPath(wf), domainName)

		var subtitle string
		if userNameElement != "" {
			subtitle = fmt.Sprintf(`Used in '%s', Group: %s`, userNameElement, domainName)
		} else {
			subtitle = fmt.Sprintf(`Used in '%s'`, domainName)
		}

		item := wf.NewItem(userNameValue).
			Subtitle(subtitle).
			Valid(true).
			Arg(userNameValue).
			Autocomplete(userNameValue).
			Copytext(userNameValue).
			Largetype(userNameValue)

		item.Cmd().Subtitle("Press Enter to paste this value directly")

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}
	}
}
