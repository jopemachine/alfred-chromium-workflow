package src

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/klauspost/lctime"
	_ "github.com/mattn/go-sqlite3"
	psl "github.com/weppos/publicsuffix-go/publicsuffix"
	// "github.com/deanishe/awgo"
)

var CheckError = func(err error) {
	if err != nil {
		panic(err)
	}
}

var EnsureDirectoryExist = func(dirPath string) {
	if !FileExist(dirPath) {
		os.Mkdir(dirPath, 0777)
	}
}

var GetDBFilePath = func(chromeProfilePath string, dbFile string) string {
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

var HandleUserQuery = func(query string) (titleQuery string, domainQuery string, isDomainSearch bool, artistQuery string, isArtistSearch bool) {
	titleQuery = ""
	domainQuery = ""
	artistQuery = ""
	isDomainSearch = false
	isArtistSearch = false

	if strings.Contains(query, "#") || strings.Contains(query, "@") {
		var words = strings.Split(query, " ")

		for _, word := range words {
			if strings.HasPrefix(word, "#") && len(word) > 1 {
				isDomainSearch = true
				domainQuery = word[1 : len(word)-1]
			} else if strings.HasPrefix(word, "@") && len(word) > 1 {
				isArtistSearch = true
				artistQuery = word[1 : len(word)-1]
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

var GetHistoryDB = func() *sql.DB {
	var targetPath = GetDBFilePath(Conf.Profile, "History")
	CopyFile(targetPath, CONSTANT.HISTORY_DB)
	db, err := sql.Open("sqlite3", CONSTANT.HISTORY_DB)
	CheckError(err)

	return db
}

var GetFaviconDB = func() *sql.DB {
	var targetPath = GetDBFilePath(Conf.Profile, "Favicons")
	CopyFile(targetPath, CONSTANT.FAVICON_DB)
	db, err := sql.Open("sqlite3", CONSTANT.FAVICON_DB)
	CheckError(err)

	return db
}

var GetWebDataDB = func() *sql.DB {
	var targetPath = GetDBFilePath(Conf.Profile, "Web Data")
	CopyFile(targetPath, CONSTANT.WEB_DATA_DB)
	db, err := sql.Open("sqlite3", CONSTANT.WEB_DATA_DB)
	CheckError(err)

	return db
}

var GetLoginDataDB = func() *sql.DB {
	var targetPath = GetDBFilePath(Conf.Profile, "Login Data")
	CopyFile(targetPath, CONSTANT.LOGIN_DATA_DB)
	db, err := sql.Open("sqlite3", CONSTANT.LOGIN_DATA_DB)
	CheckError(err)

	return db
}

var GetChromeBookmark = func() map[string]interface{} {
	var bookmarkJson map[string]interface{}
	var bookmarkFilePath = GetDBFilePath(Conf.Profile, "Bookmarks")

	bookmarkData, err := ioutil.ReadFile(bookmarkFilePath)
	CheckError(err)
	err = json.Unmarshal(bookmarkData, &bookmarkJson)
	CheckError(err)

	return bookmarkJson["roots"].(map[string]interface{})
}

// Ref: https://golangcode.com/how-to-check-if-a-string-is-a-url/
// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

var ExtractDomainName = func(url string) (domainName string) {
	if !isValidUrl(url) {
		return "unknown"
	}

	var hostname string
	if strings.Contains(url, "//") {
		hostname = strings.Split(url, "/")[2]
	} else {
		hostname = strings.Split(url, "/")[0]
	}

	hostname = strings.Split(hostname, ":")[0]
	hostname = strings.Split(hostname, "?")[0]

	domainName, err := psl.Domain(hostname)
	if err != nil {
		return hostname
	}

	return domainName
}

var ConvertChromeTimeToUnixTimestamp = func(time int64) int64 {
	return int64((math.Floor((float64(time) / 1000000) - 11644473600)))
}

var FileExist = func(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

var GetLocaleString = func(unixTime int64) string {
	return lctime.Strftime("%c", time.Unix(unixTime, 0))
}

// Used only in fetchBookmark.go of`chf` command
var ParseUserQuery = func(userQuery string) (input string, options map[string]string) {
	options = make(map[string]string)

	for _, args := range strings.Split(userQuery, " ") {
		reg, err := regexp.Compile("--[a-zA-Z\\d]*")
		CheckError(err)
		argList := strings.Split(args, " ")

		for _, arg := range argList {
			if reg.MatchString(arg) {
				key := strings.Split(strings.Split(arg, "--")[1], "=")[0]
				value := strings.Split(arg, "=")[1]
				options[key] = value
			} else {
				input += (arg + " ")
			}
		}
	}

	if strings.HasSuffix(input, " ") {
		input = strings.Trim(input, " ")
	}

	return
}

// TODO: Replace below function with stdlib's one when it is merged
// Ref: https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func StringContains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

type BookmarkItem struct {
	Id       string        `json:"id"`
	Url      string        `json:"url,omitempty"`
	Name     string        `json:"name,omitempty"`
	Children []interface{} `json:"children,omitempty"`
}

type TraverseBookmarkJsonOption struct {
	Targets []string
	Depth   int
}

// InitBookmarkJsonTraversal should be called first before calling this function
var TraverseBookmarkJSONObject func(item map[string]interface{}, options TraverseBookmarkJsonOption) []BookmarkItem
var TraverseBookmarkArray func(item []interface{}, options TraverseBookmarkJsonOption) []BookmarkItem

var InitBookmarkJsonTraversal = func() {
	TraverseBookmarkJSONObject = func(jsonObject map[string]interface{}, options TraverseBookmarkJsonOption) (result []BookmarkItem) {
		// Base case
		if options.Depth <= -1 {
			return []BookmarkItem{}
		}

		// Base case
		if jsonObject["type"] == "url" {
			if StringContains(options.Targets, "url") {
				return []BookmarkItem{
					{
						jsonObject["id"].(string),
						jsonObject["url"].(string),
						jsonObject["name"].(string),
						nil,
					},
				}
			}

			return []BookmarkItem{}
		}

		if StringContains(options.Targets, "folder") && jsonObject["type"] == "folder" {
			result = append(result, BookmarkItem{
				jsonObject["id"].(string),
				"",
				jsonObject["name"].(string),
				jsonObject["children"].([]interface{}),
			})

			childResult := TraverseBookmarkArray(jsonObject["children"].([]interface{}), options)
			result = append(result, childResult...)
			return result
		}

		target := jsonObject

		for _, child := range target {
			switch child.(type) {
			case map[string]interface{}:
				childResult := TraverseBookmarkJSONObject(child.(map[string]interface{}), options)
				result = append(result, childResult...)
			case []interface{}:
				childResult := TraverseBookmarkArray(child.([]interface{}), options)
				result = append(result, childResult...)
			}
		}

		return result
	}

	TraverseBookmarkArray = func(item []interface{}, options TraverseBookmarkJsonOption) []BookmarkItem {
		// Base case
		if options.Depth <= -1 {
			return []BookmarkItem{}
		}

		target := item
		result := []BookmarkItem{}

		for _, child := range target {
			switch child.(type) {
			case map[string]interface{}:
				childResult := TraverseBookmarkJSONObject(child.(map[string]interface{}), options)
				result = append(result, childResult...)
			case []interface{}:
				childResult := TraverseBookmarkArray(child.([]interface{}), options)
				result = append(result, childResult...)
			}
		}

		return result
	}
}
