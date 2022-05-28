package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deanishe/awgo"
)

var CacheFavicons = func(wf *aw.Workflow) {
	historyDB := GetHistoryDB(wf)
	GetFaviconDB(wf)
	faviconDBFilePath := filepath.Join(GetTempDataPath(wf), CONSTANT.FAVICON_DB)

	attachStmt, err := historyDB.Prepare(fmt.Sprintf(`ATTACH DATABASE '%s' AS favicons`, faviconDBFilePath))
	attachStmt.Exec()

	dbQuery := `
		SELECT urls.url, favicon_bitmaps.image_data, favicon_bitmaps.last_updated
			FROM urls
				LEFT OUTER JOIN icon_mapping ON icon_mapping.page_url = urls.url,
					favicon_bitmaps ON favicon_bitmaps.id =
						(SELECT id FROM favicon_bitmaps
							WHERE favicon_bitmaps.icon_id = icon_mapping.icon_id
							ORDER BY width DESC LIMIT 1)
			WHERE (urls.title LIKE '%%' OR urls.url LIKE '%%')
		`

	rows, err := historyDB.Query(dbQuery)
	CheckError(err)

	var url string
	var faviconBitmapData string
	var faviconLastUpdated string

	for rows.Next() {
		err := rows.Scan(&url, &faviconBitmapData, &faviconLastUpdated)
		CheckError(err)

		domainName := ExtractDomainName(url)
		iconPath := fmt.Sprintf(GetFaviconDirectoryPath(wf), domainName)

		if !FileExist(iconPath) {
			ioutil.WriteFile(iconPath, []byte(faviconBitmapData), os.FileMode(0777))
		}
	}

	// To send success alert
	fmt.Println(" ")
}
