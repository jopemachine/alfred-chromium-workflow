package src

import (
	"fmt"
	"github.com/deanishe/awgo"
	"strings"
)

var FetchDownloadHistory = func(wf *aw.Workflow, query string) {
	var dbQuery = `SELECT current_path, referrer, total_bytes, start_time FROM downloads ORDER BY start_time DESC`

	historyDB := GetHistoryDB()
	rows, err := historyDB.Query(dbQuery)
	CheckError(err)

	var downloadedFilePath string
	var downloadedFileFrom string
	var downloadedFileBytes int
	var downloadedStartTime int64

	for rows.Next() {
		err := rows.Scan(&downloadedFilePath, &downloadedFileFrom, &downloadedFileBytes, &downloadedStartTime)
		CheckError(err)

		fileNameArr := strings.Split(downloadedFilePath, "/")
		fileName := fileNameArr[len(fileNameArr)-1]
		domainName := ExtractDomainName(downloadedFileFrom)

		var subtitle string
		if FileExist(downloadedFilePath) {
			subtitle = "[✔]"
		} else {
			subtitle = "[✖]"
		}

		unixTimestamp := ConvertChromeTimeToUnixTimestamp(downloadedStartTime)
		localeTimeStr := GetLocaleString(unixTimestamp)

		subtitle += fmt.Sprintf(` Downloaded in %s, From '%s'`, localeTimeStr, domainName)

		item := wf.NewItem(fileName).
			Subtitle(subtitle).
			Valid(true).
			Quicklook(downloadedFilePath).
			Copytext(downloadedFilePath).
			Largetype(downloadedFilePath)

		iconPath := fmt.Sprintf(`cache/%s`, domainName)

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}
	}

	if query != "" {
		wf.Filter(query)
	}
}
