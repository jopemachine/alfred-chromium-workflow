package src

import (
	"fmt"
	"os"
	"io/ioutil"
	// "path/filepath"
	"github.com/deanishe/awgo"
)

var FetchHistory = func (wf *aw.Workflow, query string) {
	var whereStmt string
	titleQuery, domainQuery, isDomainSearch, _, _:= HandleUserQuery(query)

	if (isDomainSearch) {
		whereStmt = fmt.Sprintf(`WHERE (urls.title LIKE '%%%s%%' AND urls.url LIKE '%%%s%%')`, titleQuery, domainQuery)
	} else {
		whereStmt = fmt.Sprintf(`WHERE (urls.title LIKE '%%%s%%' OR urls.url LIKE '%%%s%%')`, titleQuery, domainQuery)
	}

	var dbQuery = fmt.Sprintf(`
		SELECT urls.id, urls.title, urls.url, urls.last_visit_time, favicon_bitmaps.image_data, favicon_bitmaps.last_updated
			FROM urls
					LEFT OUTER JOIN icon_mapping ON icon_mapping.page_url = urls.url,
							favicon_bitmaps ON favicon_bitmaps.id =
								(SELECT id FROM favicon_bitmaps
										WHERE favicon_bitmaps.icon_id = icon_mapping.icon_id
										ORDER BY width DESC LIMIT 1)
			%s
		ORDER BY last_visit_time DESC
	`, whereStmt)

	historyDB := GetHistoryDB()
	GetFaviconDB()

	attachStmt, err := historyDB.Prepare(fmt.Sprintf(`ATTACH DATABASE './%s' AS favicons`, CONSTANT.FAVICON_DB))
	attachStmt.Exec()

	CheckError(err)

	rows, err := historyDB.Query(dbQuery)
	CheckError(err)

	var urlTitle string
	var url string
	var urlId string
	var urlLastVisitTime int64

	var faviconBitmapData string
	var faviconLastUpdated string

	var itemCount = 0
	var previousTitle = ""

	EnsureDirectoryExist("./cache")

	for rows.Next() {
		if itemCount >= int(Conf.ResultLimitCount) {
			break
		}

		err := rows.Scan(&urlId, &urlTitle, &url, &urlLastVisitTime, &faviconBitmapData, &faviconLastUpdated)
		CheckError(err)

		if previousTitle == urlTitle {
			continue
		}

		domainName := ExtractDomainName(url)
		unixTimestamp := ConvertChromeTimeToUnixTimestamp(urlLastVisitTime)
		localeTimeStr := GetLocaleString(unixTimestamp)

		// iconPath, err := filepath.Abs(fmt.Sprintf(`./cache/%s.png`, domainName))
		iconPath := fmt.Sprintf(`cache/%s.png`, domainName)
		CheckError(err)

		if !FileExist(iconPath) {
			ioutil.WriteFile(iconPath, []byte(faviconBitmapData), os.FileMode(0777))
		}

		wf.NewItem(urlTitle).
			Subtitle(fmt.Sprintf(`From '%s', In '%s'`, domainName, localeTimeStr)).
			Valid(true).
			Quicklook(url).
			Var("type", "url").
			Var("url", url).
			Copytext(url).
			Largetype(url).
			Icon(&aw.Icon{iconPath, ""})

		previousTitle = urlTitle
		itemCount += 1
	}
}


