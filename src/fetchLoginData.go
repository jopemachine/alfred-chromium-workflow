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
		if userNameValue == "" {
			continue
		}

		CheckError(err)

		domainName := ExtractDomainName(originUrl)
		iconPath := fmt.Sprintf(GetFaviconDirectoryPath(wf), domainName)

		var subtitle string
		if userNameElement != "" {
			subtitle = fmt.Sprintf(`Used by '%s', Group: %s`, userNameElement, domainName)
		} else {
			subtitle = fmt.Sprintf(`Used by '%s'`, domainName)
		}

		item := wf.NewItem(userNameValue).
			Subtitle(subtitle).
			Valid(true).
			Copytext(userNameValue).
			Largetype(userNameValue)

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}
	}
}
