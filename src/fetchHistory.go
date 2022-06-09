package src

import (
	"fmt"
	"github.com/deanishe/awgo"
	"strings"
)

var FetchHistory = func(wf *aw.Workflow, query string) {
	titleQuery, domainQuery, isDomainSearch := ParseUserQuery(query)

	var dbQuery = fmt.Sprintf(`
		SELECT urls.id, urls.title, urls.url, urls.last_visit_time FROM urls
		WHERE urls.title LIKE '%%%s%%'
		ORDER BY last_visit_time DESC
	`, titleQuery)

	historyDB := GetHistoryDB(wf)

	rows, err := historyDB.Query(dbQuery)
	CheckError(err)

	var urlTitle string
	var url string
	var urlId string
	var urlLastVisitTime int64

	var itemCount = 0
	var previousTitle = ""

	for rows.Next() {
		if itemCount >= Conf.ResultCountLimit {
			break
		}

		err := rows.Scan(&urlId, &urlTitle, &url, &urlLastVisitTime)
		CheckError(err)

		if previousTitle == urlTitle {
			continue
		}

		domainName := ExtractDomainName(url)
		if isDomainSearch && !strings.Contains(domainName, domainQuery) {
			continue
		}

		unixTimestamp := ConvertChromeTimeToUnixTimestamp(urlLastVisitTime)
		localeTimeStr := GetLocaleString(unixTimestamp)

		item := wf.NewItem(urlTitle).
			Subtitle(localeTimeStr).
			Valid(true).
			Quicklook(url).
			Autocomplete(urlTitle).
			Arg(url).
			Copytext(url).
			Largetype(urlTitle)

		item.Cmd().Subtitle("Press Enter to copy this url to clipboard")

		iconPath := fmt.Sprintf(`%s/%s.png`, GetFaviconDirectoryPath(wf), domainName)

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}

		previousTitle = urlTitle
		itemCount += 1
	}
}
