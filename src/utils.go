package src

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"math"
	"errors"
	"time"
	"encoding/json"
	"github.com/klauspost/lctime"

	_ "github.com/mattn/go-sqlite3"
	psl "github.com/weppos/publicsuffix-go/publicsuffix"
	// "github.com/deanishe/awgo"
)

var CheckError = func (err error) {
	if err != nil {
		panic(err)
	}
}

var GetDBFilePath = func (chromeProfilePath string, dbFile string) string {
	var targetPath string

	switch Conf.Browser {
		case "Chrome Canary":
			targetPath = `/Users/%s/Library/Application Support/Google/Chrome Canary/%s/%s`
		case "Edge":
			targetPath = `/Users/%s/Library/Application Support/Microsoft Edge/%s/%s`
		case "Chromium":
			// 'Chrome Cloud Enrollment' could be wrong (not sure)
			targetPath = `/Users/%s/Library/Application Support/Google/Chrome Cloud Enrollment/%s/%s`
		default:
			targetPath = `/Users/%s/Library/Application Support/Google/Chrome/%s/%s`
	}

	user, err := user.Current()
	CheckError(err)
	userName := user.Username

	return fmt.Sprintf(targetPath, userName, chromeProfilePath, dbFile)
}

var HandleUserQuery = func (query string) (titleQuery string, domainQuery string, isDomainSearch bool, artistQuery string, isArtistSearch bool) {
	if strings.Contains(query, "#") || strings.Contains(query, "@") {
		var words = strings.Split(query, " ")

		for _, word := range words {
			if strings.HasPrefix(word, "#") {
				isDomainSearch = true
				domainQuery = word[1: len(word) - 1]
			} else if strings.HasPrefix(word, "@") {
				isArtistSearch = true
				artistQuery = word[1: len(word) - 1]
			} else {
				// TODO: Refactor below logic using `strings.Join`
				if titleQuery == "" {
					titleQuery += word
				} else {
					titleQuery += " " + word
				}
			}
		}
	} else {
		titleQuery = query
	}

	return
}

func CopyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		panic(err)
	}
	err = out.Sync()
	if err != nil {
		panic(err)
	}
}

var GetHistoryDB = func () (*sql.DB) {
	var targetPath = GetDBFilePath(Conf.ChromeProfile, "History")
	CopyFile(targetPath, CONSTANT.HISTORY_DB)
	db, err := sql.Open("sqlite3", CONSTANT.HISTORY_DB)
	CheckError(err)

	return db
}

var GetFaviconDB = func () (*sql.DB) {
	var targetPath = GetDBFilePath(Conf.ChromeProfile, "Favicons")
	CopyFile(targetPath, CONSTANT.FAVICON_DB)
	db, err := sql.Open("sqlite3", CONSTANT.FAVICON_DB)
	CheckError(err)

	return db
}

var GetWebDataDB = func () (*sql.DB) {
	var targetPath = GetDBFilePath(Conf.ChromeProfile, "Web Data")
	CopyFile(targetPath, CONSTANT.WEB_DATA_DB)
	db, err := sql.Open("sqlite3", CONSTANT.WEB_DATA_DB)
	CheckError(err)

	return db
}

var GetLoginDataDB = func () (*sql.DB) {
	var targetPath = GetDBFilePath(Conf.ChromeProfile, "Login Data")
	CopyFile(targetPath, CONSTANT.LOGIN_DATA_DB)
	db, err := sql.Open("sqlite3", CONSTANT.LOGIN_DATA_DB)
	CheckError(err)

	return db
}

var GetChromeBookmark = func () (bookmarkJson []map[string]interface{}) {
	var bookmarkFilePath = GetDBFilePath(Conf.ChromeProfile, "Bookmarks")

	bookmarkData, err := ioutil.ReadFile(bookmarkFilePath)
	CheckError(err)
	err = json.Unmarshal(bookmarkData, &bookmarkJson)
	CheckError(err)

	return
}

var DeleteDuplcatedItems = func (historys []map[string]string, itemLimitCount int) (result []map[string]string, deletedCount int) {
	var previousTitle string

	for idx, item := range historys {
		if idx >= itemLimitCount {
			break
		}

		if item["title"] == previousTitle {
			deletedCount += 1
		} else {
			previousTitle = item["title"]
			result = append(result, item)
		}
	}

	return
}

var ExtractDomainName = func (url string) (domainName string) {
	var hostname string
	if strings.Contains(url, "//") {
		hostname = strings.Split(url, "/")[2]
	} else {
		hostname = strings.Split(url, "/")[0]
	}

	hostname = strings.Split(hostname, ":")[0]
	hostname = strings.Split(hostname, "?")[0]

	domainName, err := psl.Domain(hostname)
	CheckError(err)

	return domainName
}

var ConvertChromeTimeToUnixTimestamp = func (time int64) int64 {
	return int64((math.Floor(((float64(time) / 1000000)) - 11644473600)) * 1000)
}

var FileExist = func (filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

var GetLocaleString = func (unixTime int64) string {
	return lctime.Strftime("%c", time.Unix(unixTime, 0))
}
