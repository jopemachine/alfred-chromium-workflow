package src

import (
	"fmt"
	"github.com/deanishe/awgo"
	"strings"
)

var FetchDownloadHistory = func(wf *aw.Workflow, query string, showOnlyExistingFiles bool) {
	var dbQuery = `SELECT current_path, referrer, total_bytes, start_time FROM downloads ORDER BY start_time DESC`

	historyDB := GetHistoryDB(wf)
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
		if fileName == "" {
			continue
		}

		var subtitle string
		if FileExist(downloadedFilePath) {
			subtitle = "[✔] "
		} else {
			if showOnlyExistingFiles == true {
				continue
			}
			subtitle = "[✖] "
		}

		domainName := ExtractDomainName(downloadedFileFrom)
		unixTimestamp := ConvertChromeTimeToUnixTimestamp(downloadedStartTime)
		localeTimeStr := GetLocaleString(unixTimestamp)

		subtitle += localeTimeStr

		item := wf.NewItem(fileName).
			Subtitle(subtitle).
			Valid(true).
			Arg(downloadedFilePath).
			Quicklook(downloadedFilePath).
			Autocomplete(fileName).
			Copytext(downloadedFilePath).
			Largetype(downloadedFilePath)

		iconPath := fmt.Sprintf(`%s/%s.png`, GetFaviconDirectoryPath(wf), domainName)

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}
	}

	if query != "" {
		wf.Filter(query)
	}
}
