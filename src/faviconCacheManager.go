package src

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/deanishe/awgo"
)

var IsFaviconCacheExpired = func(wf *aw.Workflow) bool {
	faviconCacheDir := GetFaviconDirectoryPath(wf)
	fileInfo, err := os.Stat(faviconCacheDir)
	CheckError(err)

	// For test
	// return time.Since(fileInfo.ModTime()) > time.Minute

	// Update favicon every three days
	return time.Since(fileInfo.ModTime()) > time.Hour*24*3
}

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
		iconPath := fmt.Sprintf(`%s/%s.png`, GetFaviconDirectoryPath(wf), domainName)

		if !FileExist(iconPath) {
			ioutil.WriteFile(iconPath, []byte(faviconBitmapData), os.FileMode(0777))
		}
	}

	// Change folder's atime and mtime both to currenttime even if nothing changed
	currenttime := time.Now().Local()
	err = os.Chtimes(GetFaviconDirectoryPath(wf), currenttime, currenttime)
	CheckError(err)

	// To send success alert
	fmt.Println(" ")
}

var EnsureFaviconCacheUptodated = func(wf *aw.Workflow) {
	faviconCacheDir := GetFaviconDirectoryPath(wf)

	// To avoid refreshing cache delay result, run the task in background when update the favicons
	if isEmpty, err := IsEmptyDirectory(faviconCacheDir); isEmpty || err != nil {
		defer func() {
			err := recover()
			log.Println("Error occurs in caching favicon: ", err)
		}()

		CacheFavicons(wf)
	} else if IsFaviconCacheExpired(wf) {
		cmd := exec.Command(os.Args[0], "cache-favicons")
		wf.RunInBackground("favicon-cache", cmd)
	}
}
